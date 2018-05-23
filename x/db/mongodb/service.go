package mongodb

import "gopkg.in/mgo.v2"

type Service struct {
	baseSession *mgo.Session
	queue       chan int
	URL         string
	DbUser      string
	DbPass      string
	Open        int
}

var service Service

func (s *Service) New() error {
	var err error
	var maxPool = MongoConfig.MaxPool
	s.queue = make(chan int, maxPool)
	for i := 0; i < maxPool; i = i + 1 {
		s.queue <- 1
	}
	s.Open = 0
	var dialInfo = &mgo.DialInfo{
		Addrs:    []string{s.URL},
		Username: s.DbUser,
		Password: s.DbPass,
	}
	s.baseSession, err = mgo.DialWithInfo(dialInfo)
	return err
}

func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}
