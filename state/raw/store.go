package raw

// Store provides primitives for storing, accessing and manipulating swarm objects.
type Store interface {
	CreateObject(key string, o interface{}) error
	DeleteObject(key string) error
	Object(key string) interface{}
	ListObjects(key string) []interface{}
}
