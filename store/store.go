package store

import (
	"errors"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Strategy int

const (
	Absolute     Strategy      = 0
	Sliding      Strategy      = 1
	NoExpiration time.Duration = 0
)

type item struct {
	data               interface{}
	expireAt           int64
	expiration         time.Duration
	expirationStrategy Strategy
}

type metrics struct {
	putCount prometheus.Counter
	getCount prometheus.Counter
	delCount prometheus.Counter
}

type Store struct {
	items                 sync.Map
	counters              map[string]int64
	defaultExpiration     time.Duration
	cleanupInterval       time.Duration
	defaultExpireStrategy Strategy
	metrics               *metrics
	mutex                 *sync.Mutex
}

func New(defaultExpiration time.Duration,
	defaultExpireStrategy Strategy,
	cleanupInterval time.Duration) *Store {
	s := &Store{
		counters:              make(map[string]int64),
		defaultExpiration:     defaultExpiration,
		cleanupInterval:       cleanupInterval,
		defaultExpireStrategy: defaultExpireStrategy,
		metrics: &metrics{
			putCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "total_put",
				Help: "Total number of put requests",
			}),
			getCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "total_get",
				Help: "Total number of get requests",
			}),
			delCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "total_del",
				Help: "Total number of del requests",
			}),
		},
		mutex: &sync.Mutex{},
	}

	prometheus.MustRegister(s.metrics.putCount, s.metrics.getCount, s.metrics.delCount)

	go s.startEviction()
	return s
}

func (s *Store) startEviction() {
	ticker := time.NewTicker(s.cleanupInterval)

	for t := range ticker.C {
		now := t.UnixNano()
		s.items.Range(func(key, val interface{}) bool {
			value := val.(*item)
			if value.expireAt > 0 && now > value.expireAt {
				s.items.Delete(key)
			}
			return true
		})
	}
}

func (s *Store) Put(key string, value interface{}) {
	s.PutWithExpire(key, value, s.defaultExpiration, s.defaultExpireStrategy)
}

func (s *Store) PutWithoutExpire(key string, value interface{}) {
	s.PutWithExpire(key, value, NoExpiration, s.defaultExpireStrategy)
}

func (s *Store) PutWithExpire(key string, value interface{}, expiration time.Duration, strategy Strategy) {
	var expireAt int64
	if expiration.Nanoseconds() > 0 {
		expireAt = time.Now().Add(expiration).UnixNano()
	}

	v := &item{
		data:               value,
		expireAt:           expireAt,
		expiration:         expiration,
		expirationStrategy: strategy,
	}

	s.items.Store(key, v)
	go func() {
		s.metrics.putCount.Inc()
	}()
}

func (s *Store) Get(key string) (interface{}, bool) {
	go func() {
		s.metrics.getCount.Inc()
	}()

	value, found := s.items.Load(key)
	if !found {
		return nil, false
	} else {
		value := value.(*item)
		if value.expireAt == 0 { // No expiration
			return value.data, true
		}

		now := time.Now().UnixNano()

		if now > value.expireAt {
			return nil, false
		}

		if value.expirationStrategy == Sliding {
			value.expireAt += value.expiration.Nanoseconds()
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
	go func() {
		s.metrics.delCount.Inc()
	}()

	s.items.Delete(key)
}

func (s *Store) IncCounter(key string, valueToAdd int64) int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	counter, found := s.counters[key]
	if !found {
		s.counters[key] = valueToAdd
		return valueToAdd
	}
	counter += valueToAdd
	s.counters[key] = counter
	return counter
}

func (s *Store) DecCounter(key string, valueToDec int64) int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	counter, found := s.counters[key]
	if !found {
		s.counters[key] = valueToDec
		return valueToDec
	}
	counter -= valueToDec
	s.counters[key] = counter
	return counter
}

func (s *Store) GetCounter(key string) (int64, error) {
	s.mutex.Lock()
	counter, found := s.counters[key]
	s.mutex.Unlock()

	if !found {
		return 0, errors.New("Counter not found")
	}
	return counter, nil
}
