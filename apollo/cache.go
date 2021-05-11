package apollo

import (
	"encoding/json"
	"os"

	"github.com/zhucenichenghao/configure/rlog"
)

func (c *Client) dumpToFile() error {
	f, err := os.OpenFile(c.cacheFileName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	if err != nil {
		rlog.Error("OpenFile failed", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return json.NewEncoder(f).Encode(&c.Configs)
}

func (c *Client) loadFromCache() error {
	f, err := os.OpenFile(c.cacheFileName(), os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	if err := json.NewDecoder(f).Decode(&c.Configs); err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheFileName() string {
	return "./apollo_configs_" + c.Conf.AppID
}
