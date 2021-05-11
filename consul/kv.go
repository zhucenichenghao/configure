package consul

import (
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/zhucenichenghao/configure/rlog"
)

func compareMap(v1 map[string]string, v2 map[string]string) bool {
	if len(v1) != len(v2) {
		return false
	}
	for k, s1 := range v1 {
		s2, ok := v2[k]
		if !ok {
			return false
		}
		if s1 != s2 {
			return false
		}
	}
	return true
}

func deepCopy(v map[string]string) map[string]string {
	dup := make(map[string]string, len(v))
	for k, s := range v {
		dup[k] = s
	}
	return dup
}

// WatchKV 对KV的watch
func WatchKV(client *api.Client, path string, config chan map[string]string) {
	var lastIndex uint64
	var lastValue map[string]string
	for {
		value, index, err := ListKVDict(client, path, lastIndex)
		if err != nil {
			rlog.Error("consul configure: listkv failed", map[string]interface{}{
				"error": err.Error(),
				"path":  path,
			})
			time.Sleep(time.Second)
			continue
		}
		if index != lastIndex {
			lastIndex = index
		}
		if !compareMap(lastValue, value) {
			lastValue = value
			config <- deepCopy(value)
		}
	}
}

// ListKVDict 以map的形式列出path下kv
func ListKVDict(client *api.Client, path string, waitIndex uint64) (map[string]string, uint64, error) {
	q := &api.QueryOptions{RequireConsistent: true, WaitIndex: waitIndex}
	kvpairs, meta, err := client.KV().List(path, q)
	if err != nil {
		return nil, 0, err
	}
	dict := make(map[string]string)
	if len(kvpairs) == 0 {
		return dict, meta.LastIndex, nil
	}
	for _, kvpair := range kvpairs {
		val := strings.TrimSpace(string(kvpair.Value))
		dict[kvpair.Key] = val
	}
	return dict, meta.LastIndex, nil
}

// GetKV 获取key对应的value
func GetKV(client *api.Client, key string, waitIndex uint64) (string, uint64, error) {
	q := &api.QueryOptions{RequireConsistent: true, WaitIndex: waitIndex}
	kvpair, meta, err := client.KV().Get(key, q)
	if err != nil {
		return "", 0, err
	}
	if kvpair == nil {
		return "", meta.LastIndex, nil
	}
	return strings.TrimSpace(string(kvpair.Value)), meta.LastIndex, nil
}

// PutKV 写key对应的value值
func PutKV(client *api.Client, key, value string, index uint64) (bool, error) {
	p := &api.KVPair{Key: key, Value: []byte(value), ModifyIndex: index}
	var err error = nil
	ok := true
	if index == 0 {
		_, err = client.KV().Put(p, nil)
	} else {
		ok, _, err = client.KV().CAS(p, nil)
	}
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("putkv failed")
	}
	return true, nil
}

// DelKV 删除key
func DelKV(client *api.Client, key string, index uint64) (bool, error) {
	p := &api.KVPair{Key: key, ModifyIndex: index}
	var err error = nil
	ok := true
	if index == 0 {
		_, err = client.KV().Delete(key, nil)
	} else {
		ok, _, err = client.KV().DeleteCAS(p, nil)
	}
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("delkv failed")
	}
	return true, nil
}
