package mongodb

import (
	"cetm_booking/x/rest"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERR_EXIST = rest.BadRequest("exists unique in db")
)

type Table struct {
	*mgo.Collection
	Name   string
	Prefix string
	Length int
}

func NewTable(name, prefix string, length int) *Table {
	fmt.Println("DB: " + name)
	return &Table{
		Collection: NewCollection(name),
		Name:       name,
		Prefix:     prefix,
		Length:     length,
	}
}

func (t *Table) Create(model IModel) error {
	model.BeforeCreate(t.Prefix, t.Length)
	var err = t.Insert(model)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}

func (t *Table) CreateUnique(query bson.M, model IModel) error {
	count, err := t.CountWhere(query)
	if err == nil {
		if count == 0 {
			return t.Create(model)
		}
		logDB.Errorln(ERR_EXIST)
		return ERR_EXIST
	}
	logDB.Errorln(err)
	return err
}

func (t *Table) CountWhere(query bson.M) (int, error) {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var c, err = t.Find(query).Count()
	if err != nil {
		logDB.Errorln(err)
	}
	return c, err
}

func (t *Table) FindWhere(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var err = t.Find(query).All(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}
func (t *Table) FindOne(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	var err = t.Find(query).One(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}
func (t *Table) FindByID(id string, result interface{}) error {
	var err = t.FindId(id).One(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}

func (t *Table) DeleteByID(id string) error {
	var err = t.UpdateId(id, bson.M{"$set": bson.M{"updated_at": 0}})
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}

func (t *Table) UnsafeUpdateByID(id string, data interface{}) error {
	var err = t.UpdateId(id, bson.M{"$set": data})
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}

func (t *Table) UnsafeFindSort(queryMatch bson.M, fields string, result interface{}) error {
	var err = t.Find(queryMatch).Sort(fields).All(result)
	if err != nil {
		logDB.Errorln(err)
	}
	return err
}
