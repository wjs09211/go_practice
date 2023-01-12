package localcache

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CacheSuite struct {
	suite.Suite
	cache *localCache
}

func (s *CacheSuite) SetupSuite() {
	// call first, only once
}

func (s *CacheSuite) SetupTest() {
	// call this before each test case
	s.cache = New().(*localCache)
}

func (s *CacheSuite) TestSet() {
	tests := []struct {
		Desc  string
		Key   string
		Value any
		Error error
	}{
		{
			Desc:  "set int",
			Key:   "int",
			Value: 1,
			Error: nil,
		},
		{
			Desc:  "set string",
			Key:   "string",
			Value: "string",
			Error: nil,
		},
		{
			Desc:  "set array",
			Key:   "array",
			Value: []int{1, 2, 3},
			Error: nil,
		},
	}
	for _, t := range tests {
		err := s.cache.Set(t.Key, t.Value, time.Minute)
		s.Require().Equal(err, t.Error, t.Desc)
		if err == nil {
			s.Require().Equal(s.cache.hash[t.Key].data, t.Value, t.Desc)
		}
	}

}

func (s *CacheSuite) TestGet() {
	testKey := "key"
	testValue := "value"
	desc := "normal get"
	s.cache.hash[testKey] = &cacheData{testValue, time.Now().Add(time.Minute)}
	val, err := s.cache.Get(testKey)
	s.Require().Equal(err, nil, desc)
	s.Require().Equal(val, testValue, desc)
}

func (s *CacheSuite) TestGetNoExistKey() {
	testKey := "key"
	desc := "no exist key"
	val, err := s.cache.Get(testKey)
	s.Require().Equal(err, ErrKeyNotExist, desc)
	s.Require().Equal(val, nil, desc)
}

func (s *CacheSuite) TestCacheExpired() {
	testKey := "key"
	testValue := "value"
	desc := "cache expired should return nil value"
	_ = s.cache.Set(testKey, testValue, time.Millisecond*50)
	time.Sleep(time.Millisecond * 100)
	val, err := s.cache.Get(testKey)
	s.Require().Equal(err, ErrKeyNotExist, desc)
	s.Require().Equal(val, nil, desc)
}

func TestStart(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
