package mongodb

import (
	"cetm_booking/x/mlog"
	"gopkg.in/mgo.v2"
)

var logDB = mlog.NewTagLog("MONGO_DB")

type Database struct {
	s       *mgo.Session
	name    string
	session *mgo.Database
}

func (db *Database) Connect() {

	db.s = service.Session()
	session := *db.s.DB(db.name)
	db.session = &session

}

func newDBSession(name string) *Database {

	var db = Database{
		name: name,
	}
	db.Connect()
	return &db
}
