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

func (s *Store) GetInt(key string) (int, bool) {
	value, found := s.Get(key)
	if found {
		return value.(int), found
	}

	return 0, false
}

func (s *Store) GetInt16(key string) (int16, bool) {
	value, found := s.Get(key)
	if found {
		return value.(int16), found
	}

	return 0, false
}

func (s *Store) GetUint16(key string) (uint16, bool) {
	value, found := s.Get(key)
	if found {
		return value.(uint16), found
	}

	return 0, false
}

func (s *Store) GetInt32(key string) (int32, bool) {
	value, found := s.Get(key)
	if found {
		return value.(int32), found
	}

	return 0, false
}

func (s *Store) GetUint32(key string) (uint32, bool) {
	value, found := s.Get(key)
	if found {
		return value.(uint32), found
	}

	return 0, false
}

func (s *Store) GetInt64(key string) (int64, bool) {
	value, found := s.Get(key)
	if found {
		return value.(int64), found
	}

	return 0, false
}

func (s *Store) GetUint64(key string) (uint64, bool) {
	value, found := s.Get(key)
	if found {
		return value.(uint64), found
	}

	return 0, false
}

func (s *Store) GetBool(key string) (bool, bool) {
	value, found := s.Get(key)
	if found {
		return value.(bool), found
	}

	return false, false
}

func (s *Store) GetString(key string) (string, bool) {
	value, found := s.Get(key)
	if found {
		return value.(string), found
	}

	return "", false
}

func (s *Store) GetByteArray(key string) ([]byte, bool) {
	value, found := s.Get(key)
	if found {
		return value.([]byte), found
	}

	return make([]byte, 0), false
}

func (s *Store) Del(key string) {
	s.items.Delete(key)
}
