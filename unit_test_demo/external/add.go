package external

type cache struct{}

type Cache interface {
	Save(input interface{}) error
}

func NewCacheClient() Cache {
	return &cache{}
}
func (m *cache) Save(input interface{}) error {
	return nil
}

//func Save(input interface{}) error {
//	return nil
//}
