package store

import (
	"testing"
	"time"
)

func TestPutAndGet(t *testing.T) {
	kv := New(1*time.Second, Absolute, 30*time.Second)
	kv.Put("SomeKey", 765)
	val, found := kv.GetInt("SomeKey")
	if !found {
		t.Error("Key must be exist")
	}
	if val != 765 {
		t.Error("Invalid value")
	}

	time.Sleep(2 * time.Second)
	_, found = kv.Get("SomeKey")
	if found {
		t.Errorf("Key '%s' has not been expired", "Key2")
	}
}

func TestDel(t *testing.T) {
	kv := New(1*time.Second, Absolute, 30*time.Second)
	kv.PutWithoutExpire("Key2", "Hello")
	kv.Del("Key2")
	_, found := kv.Get("Key2")
	if found {
		t.Errorf("Key '%s' has not been expired", "Key2")
	}
}

func TestCounter(t *testing.T) {
	kv := New(1*time.Second, Absolute, 30*time.Second)
	kv.IncCounter("online_users", 1)
	kv.IncCounter("online_users", 1)
	kv.IncCounter("online_users", 1)
	kv.IncCounter("online_users", 1)
	kv.DecCounter("online_users", 1)
	onlineUsers, _ := kv.GetCounter("online_users")
	if onlineUsers != 3 {
		t.Error("Invalid countet value")
	}
}
