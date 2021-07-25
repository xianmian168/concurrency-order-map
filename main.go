package main

import (
	"container/list"
	"log"
	"sync"
)

type Map struct {
	m map[string]string
	o *list.List
	l sync.RWMutex
}

func New() *Map {
	return &Map{
		m: make(map[string]string),
		o: list.New(),
	}
}

func (m *Map) Put(key string, value string) {
	m.l.Lock()
	defer m.l.Unlock()
	_, found := m.m[key]
	if found {
		return
	}
	m.o.PushFront(key)
	m.m[key] = value
}

func (m *Map) Get(key string) (value string, found bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	value, found = m.m[key]
	return
}

func (m *Map) Remove(key string) {
	m.l.Lock()
	defer m.l.Unlock()
	_, found := m.m[key]
	if !found {
		return
	}
	delete(m.m, key)
	for e := m.o.Front(); e != nil; e = e.Next() {
		if e.Value == key {
			m.o.Remove(e)
		}
	}
}

func (m *Map) Empty() bool {
	return m.Size() == 0
}

func (m *Map) Size() int {
	return m.o.Len()
}

func (m *Map) Keys() []string {
	keys := make([]string, 0, len(m.m))
	for e := m.o.Back(); e != nil; e = e.Prev() {
		v, found := e.Value.(string)
		if !found {
			continue
		}
		keys = append(keys, v)
	}
	return keys
}

func (m *Map) Values() []string {
	values := make([]string, 0, m.Size())
	for e := m.o.Back(); e != nil; e = e.Prev() {
		k, found := e.Value.(string)
		if !found {
			continue
		}
		v, found := m.Get(k)
		if !found {
			continue
		}
		values = append(values, v)
	}
	return values
}

func main() {
	m := New()
	m.Put("2", "1")
	m.Put("1", "2")
	v, found := m.Get("1")
	log.Println(v, found)
	// 2 true
	log.Printf("%+v\n", m.Values())
	// [1 2]
	log.Printf("%+v\n", m.Keys())
	// [2 1]
}
