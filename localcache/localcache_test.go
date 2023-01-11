package localcache

import (
	"fmt"
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
	fmt.Println("SetupTest")
	s.cache = New().(*localCache)
}

func (s *CacheSuite) TestSet() {
	testKey := "key"
	testValue := "value"
	desc := "normal set"
	err := s.cache.Set(testKey, testValue, time.Minute)
	s.Require().Equal(err, nil, desc)
	s.Require().Equal(s.cache.hash[testKey].data, testValue, desc)

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
	desc := "cache expired should return nil"
	_ = s.cache.Set(testKey, testValue, time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	val, err := s.cache.Get(testKey)
	s.Require().Equal(err, nil, desc)
	s.Require().Equal(val, nil, desc)
}

func TestStart(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
