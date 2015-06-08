package data

// Storer implements a Get method
type Storer interface {
	Store(ns, key string, value interface{}) error
}

// Getter implements a Get method
type Getter interface {
	Get(ns, key string) (interface{}, bool)
}

// StoreGetter implements both the Storer and Getter interfaces
type StoreGetter interface {
	Storer
	Getter
}

// compile-time check we're still implementing the interface
var _ StoreGetter = &inMemory{}

// Storage holds the storage backend
var Storage = StoreGetter(&inMemory{data: make(map[string]interface{})})

type inMemory struct {
	data map[string]interface{}
}

func (i *inMemory) Get(ns, key string) (interface{}, bool) {
	v, ok := i.data[ns+"."+key]
	return v, ok
}

func (i *inMemory) Store(ns, key string, value interface{}) error {
	i.data[ns+"."+key] = value
	return nil
}
