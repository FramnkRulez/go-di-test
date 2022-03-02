package main

import (
	"fmt"
	"reflect"

	"github.com/FramnkRulez/go-di-test/pkg/di"
)

// define logger interface
type logger interface {
	Log(msg string)
}

// define logger implementation
type loggerImpl struct {
}

func (l loggerImpl) Log(msg string) {
	fmt.Println(msg)
}

func (l loggerImpl) Init() {

}

// define keyvaluestore interface
type keyvaluestore interface {
	AddValue(key string, value string)
	GetValue(key string) string
}

// define in-memory keyvalue store implementation
type InMemoryKeyValueStore struct {
	l      logger
	kvpMap map[string]string
}

func (kvp *InMemoryKeyValueStore) Init(l logger) {
	kvp.l = l
	kvp.kvpMap = make(map[string]string)
}

func (kvp InMemoryKeyValueStore) AddValue(key string, value string) {
	kvp.l.Log(fmt.Sprintf("AddValue called with key %v and value %v", key, value))
	kvp.kvpMap[key] = value

	fmt.Printf("Key value map now has %v values\n", len(kvp.kvpMap))
}

func (kvp InMemoryKeyValueStore) GetValue(key string) string {
	return kvp.kvpMap[key]
}

// container resolve helper
func resolveKeyValueStore(c di.Container) (keyvaluestore, error) {
	kvs, err := c.Resolve(reflect.TypeOf((*keyvaluestore)(nil)).Elem())

	var k keyvaluestore
	k = kvs.(keyvaluestore)
	return k, err
}

func main() {
	fmt.Printf("Go Dependency Injection Test\n")

	// create our container and register our dependencies with it
	var c di.Container

	// register the logger interface (implemented by loggerImpl)
	c.Register(reflect.TypeOf((*logger)(nil)).Elem(), reflect.TypeOf(loggerImpl{}))

	// register the keyvaluestore interface (implemented by InMemoryKeyValueStore, depends on logger)
	c.Register(reflect.TypeOf((*keyvaluestore)(nil)).Elem(), reflect.TypeOf(InMemoryKeyValueStore{}))

	// Method #1 (no wrapper) - resolve keyvaluestore from dependencies by force casting return
	kvs, _ := c.Resolve(reflect.TypeOf((*keyvaluestore)(nil)).Elem())
	var store keyvaluestore = kvs.(keyvaluestore)
	store.AddValue("one", "two")

	// Method #2 (wrapper) - resolve keyvaluestore without casting using a wrapper func
	store2, _ := resolveKeyValueStore(c)
	store2.AddValue("three", "four")
}
