package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/zhucenichenghao/configure"
	"github.com/zhucenichenghao/configure/rlog"
)

type client struct {
	consulCli *api.Client // consul连接
}

// NewConsulClient 创建基于consul的config client
func NewConsulClient(consulAddr string) configure.Client {
	c, err := api.NewClient(&api.Config{Address: consulAddr, Scheme: "http"})
	if err != nil {
		rlog.Error("configure client: new consul api failed ", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}
	cli := &client{
		consulCli: c,
	}
	return cli
}

// GetKeyValue ...
func (c *client) GetKeyValue(path *configure.Path) (map[string]string, error) {
	value, _, err := ListKVDict(c.consulCli, path.Key, 0)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// WatchKeyValue ...
func (c *client) WatchKeyValue(path *configure.Path) chan map[string]string {
	kv := make(chan map[string]string)
	go WatchKV(c.consulCli, path.Key, kv)
	return kv
}

// PutKeyValue ...
func (c *client) PutKeyValue(path *configure.Path, value string) error {
	_, err := PutKV(c.consulCli, path.Key, value, path.Index)
	return err
}

// DelKeyValue ...
func (c *client) DelKeyValue(path *configure.Path) error {
	_, err := DelKV(c.consulCli, path.Key, path.Index)
	return err
}
