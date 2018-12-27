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
	data               interface{}
	expireAt           time.Time
	expiration         time.Duration
	expirationStrategy Strategy
}

type Store struct {
	items                 sync.Map
	defaultExpiration     time.Duration
	cleanupInterval       time.Duration
	defaultExpireStrategy Strategy
}

func New(defaultExpiration time.Duration,
	defaultExpireStrategy Strategy,
	cleanupInterval time.Duration) *Store {
	s := &Store{
		defaultExpiration:     defaultExpiration,
		cleanupInterval:       cleanupInterval,
		defaultExpireStrategy: defaultExpireStrategy,
	}

	go s.startEviction()
	return s
}

func (s *Store) startEviction() {
	ticker := time.NewTicker(s.cleanupInterval)

	for t := range ticker.C {
		now := t.UnixNano()
		s.items.Range(func(key, val interface{}) bool {
			value := val.(*item)
			if now > value.expireAt.UnixNano() {
				s.items.Delete(key)
			}
			return true
		})
	}
}

func (s *Store) Put(key string, value interface{}) {
	s.PutWithExpire(key, value, s.defaultExpiration, s.defaultExpireStrategy)
}

func (s *Store) PutWithExpire(key string, value interface{}, expiration time.Duration, strategy Strategy) {
	v := &item{
		data:               value,
		expireAt:           time.Now().Add(expiration),
		expiration:         expiration,
		expirationStrategy: strategy,
	}

	s.items.Store(key, v)
}

func (s *Store) Get(key string) (interface{}, bool) {
	value, found := s.items.Load(key)
	if !found {
		return nil, false
	} else {
		value := value.(*item)
		now := time.Now().UnixNano()

		if now > value.expireAt.UnixNano() {
			return nil, false
		}

		if value.expirationStrategy == Sliding {
			value.expireAt = value.expireAt.Add(value.expiration)
		}

		return value.data, true
	}
}

func (s *Store) Del(key string) {
	s.items.Delete(key)
}
