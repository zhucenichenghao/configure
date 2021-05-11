package apollo

import (
	"fmt"
	"testing"

	"github.com/zhucenichenghao/configure"
	"github.com/zhucenichenghao/configure/rlog"
)

func TestNewWriteableApolloClient(t *testing.T) {
	appid := "141"
	cli := NewWriteableApolloClient(
		&Conf{
			AppID:         appid,
			Cluster:       "default",
			NameSpaceName: []string{"application"},
			IP:            "127.0.0.1:8080",
		},
		&ProtalConf{
			IP:          "127.0.0.1:8070",
			Env:         "DEV",
			AppID:       appid,
			ClusterName: "default",
			UserID:      "operate",
			Token:       "c955889febd9577b0568314feb9602f67f2f68a4",
		},
	)
	testk1 := &configure.Path{appid, "application", "testK1", 0}
	if r, err := cli.GetKeyValue(testk1); err != nil {
		rlog.Error("get kv fail", map[string]interface{}{
			"error": err.Error(),
			"appid": appid,
			"k":     testk1,
		})
		return
	} else {
		rlog.Info("get kv successful", map[string]interface{}{
			"appid": appid,
			"k":     testk1,
			"v":     r,
		})
	}

	if err := cli.PutKeyValue(testk1, "testValue"); err != nil {
		rlog.Error("put kv fail", map[string]interface{}{
			"error": err.Error(),
			"appid": appid,
			"k":     testk1,
			"v":     "testValue",
		})
		return
	}
	if r, err := cli.GetKeyValue(testk1); err != nil {
		rlog.Error("get kv fail", map[string]interface{}{
			"error": err.Error(),
			"appid": appid,
			"k":     testk1,
		})
		return
	} else {
		rlog.Info("get kv successful", map[string]interface{}{
			"appid": appid,
			"k":     testk1,
			"v":     r,
		})
	}

	if err := cli.DelKeyValue(testk1); err != nil {
		rlog.Error("del kv fail", map[string]interface{}{
			"error": err.Error(),
			"appid": appid,
			"k":     testk1,
		})
		return
	}
	if r, err := cli.GetKeyValue(testk1); err != nil {
		rlog.Error("get kv fail", map[string]interface{}{
			"error": err.Error(),
			"appid": appid,
			"k":     testk1,
		})
		return
	} else {
		rlog.Info("get kv successful", map[string]interface{}{
			"appid": appid,
			"k":     testk1,
			"v":     r,
		})
	}

	for iloop := 0; iloop < 20; iloop++ {
		go func() {
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("key%d", i)
				// if err := cli.DelKeyValue(&configure.Path{appid, "application", key, 0}); err != nil {
				// 	log.Error().Str("appid", appid).Str("key", key).Err(err).Msg("del fail")
				// }
				v := fmt.Sprintf("value%d", i)
				if err := cli.PutKeyValue(&configure.Path{appid, "application", key, 0}, v); err != nil {
					rlog.Error("put kv fail", map[string]interface{}{
						"error": err.Error(),
						"appid": appid,
					})
					return
				}
				if r, err := cli.GetKeyValue(&configure.Path{appid, "application", key, 0}); err != nil {
					rlog.Error("get kv fail", map[string]interface{}{
						"error": err.Error(),
						"appid": appid,
					})
					return
				} else {
					if gv, ok := r[key]; !ok {
						rlog.Error("k not found", map[string]interface{}{
							"appid": appid,
							"key":   key,
						})
					} else if gv != v {
						rlog.Error("v inequality", map[string]interface{}{
							"appid": appid,
							"key":   key,
							"v":     v,
							"gv":    gv,
						})
					}
				}
			}
		}()
	}

	for v := range cli.WatchKeyValue(&configure.Path{appid, "application", "config", 0}) {
		rlog.Info("kv changed", map[string]interface{}{
			"appid": appid,
			"kv":    v,
		})
	}
}
