package cache_test

import (
	"testing"
	"time"

	"github.com/umbe77/dukes/cache"
	"github.com/umbe77/dukes/datatypes"
)

func TestSet(t *testing.T) {
	c := cache.NewCache()
	var (
		err  error
		gErr error
		v    *cache.CacheValue
	)
	if err = c.Set("Key1", &cache.CacheValue{
		Kind:  datatypes.Int,
		Value: int32(10),
	}); err != nil {
		t.Errorf("Problem setting Key1")
	}
	if err = c.Set("Key2", &cache.CacheValue{
		Kind:  datatypes.String,
		Value: "Value1",
	}); err != nil {
		t.Errorf("Problem setting Key2")

	}
	key3Val := time.Now()
	if err = c.Set("Key3", &cache.CacheValue{
		Kind:  datatypes.Date,
		Value: key3Val,
	}); err != nil {
		t.Errorf("Problem setting Key3")
	}
	if err = c.Set("Key4", &cache.CacheValue{
		Kind:  datatypes.Bool,
		Value: true,
	}); err != nil {
		t.Errorf("Problem setting Key4")
	}

	v, gErr = c.Get("Key1")
	if gErr != nil {
		t.Errorf("Problem getting Key1")
	}
	if v.Kind != datatypes.Int {
		t.Errorf("Key 1 should by Int, got %v", v.Kind)
	}
	if int32(v.Value.(int32)) != 10 {
		t.Errorf("Key1 v expcetd 10, got %v", v)
	}

	v, gErr = c.Get("Key2")
	if gErr != nil {
		t.Errorf("Problem getting Key2")
	}
	if v.Kind != datatypes.String {
		t.Errorf("Key 2 should by string, got %v", v.Kind)
	}
	if v.Value.(string) != "Value1" {
		t.Errorf("Key2 v expcetd \"Value1\", got %v", v)
	}

	v, gErr = c.Get("Key3")
	if gErr != nil {
		t.Errorf("Problem getting Key3")
	}
	if v.Kind != datatypes.Date {
		t.Errorf("Key 3 should by date, got %v", v.Kind)
	}
	if v.Value.(time.Time) != key3Val {
		t.Errorf("Key3 v expcetd %v, got %v", key3Val, v)
	}

	v, gErr = c.Get("Key4")
	if gErr != nil {
		t.Errorf("Problem getting Key4")
	}
	if v.Kind != datatypes.Bool {
		t.Errorf("Key 4 should by bool, got %v", v.Kind)
	}
	if v.Value.(bool) != true {
		t.Errorf("Key4 v expcetd true, got %v", v)
	}
}

func TestDel(t *testing.T) {
	c := cache.NewCache()
	c.Set("Key1", &cache.CacheValue{
		Kind:  datatypes.Int,
		Value: 1,
	})

	if err := c.Del("Key1"); err != nil {
		t.Errorf("Error on deleting Key1: %v", err)
	}

	if v, err := c.Get("Key1"); err == nil {
		t.Errorf("Key1 still present in cache with value %v", v)
	}

}

func TestHas(t *testing.T) {
	c := cache.NewCache()
	c.Set("key1", &cache.CacheValue{
		Kind:  datatypes.String,
		Value: 1,
	})

	if !c.Has("key1") {
		t.Errorf("key1 not found")
	}

	if c.Has("key2") {
		t.Errorf("key2 found")
	}

}
