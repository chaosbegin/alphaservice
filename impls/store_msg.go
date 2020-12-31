package impls

type StoreMsg struct {
	Series string
	Host   string
	Tags   []string
	Data   []map[string]interface{}
}
