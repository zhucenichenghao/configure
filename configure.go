package configure

// Path 路径信息
type Path struct {
	AppID     string
	NameSpace string
	Key       string
	Index     uint64 // CAS
	// 0--精准匹配 1--模糊匹配
	Match int8
}

// Client 这里抽象配置中心
type Client interface {

	// GetKeyValue ...
	GetKeyValue(*Path) (map[string]string, error)

	// WatchKeyValue ...
	WatchKeyValue(*Path) chan map[string]string

	// PutKeyValue ...
	PutKeyValue(*Path, string) error

	// DelKeyValue ...
	DelKeyValue(*Path) error
}
