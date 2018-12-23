package store

import (
	"sync"
	"time"
)

type Strategy int

const (
	Absolute Strategy = 0
	Sliding  Strategy = 1
)

type item struct {
	data       interface{}
	expireAt   time.Time
	expiration time.Duration
}

type store struct {
	items             sync.Map
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	expireStrategy    Strategy
}

func New(defaultExpiration time.Duration, expireStrategy Strategy, cleanupInterval time.Duration) *store {
	return &store{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		expireStrategy:    expireStrategy,
	}
}

func (s *store) Put(key string, value interface{}) {
	s.PutWithExpire(key, value, s.defaultExpiration)
}

func (s *store) PutWithExpire(key string, value interface{}, expiration time.Duration) {
	v := &item{
		data:       value,
		expireAt:   time.Now().Add(expiration),
		expiration: expiration,
	}

	s.items.Store(key, v)
}

func (s *store) Get(key string) (interface{}, bool) {
	value, found := s.items.Load(key)
	if !found {
		return nil, false
	} else {
		value := value.(*item)
		now := time.Now().UnixNano()

		if now > value.expireAt.UnixNano() {
			return nil, false
		}

		if s.expireStrategy == Sliding {
			value.expireAt = value.expireAt.Add(value.expiration)
		}

		return value.data, true
	}
}

func (s *store) Del(key string) {
	s.items.Delete(key)
}
