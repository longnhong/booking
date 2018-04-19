package notify

import (
	"cetm_booking/x/rest"
	"gopkg.in/mgo.v2/bson"
)

func (noti *Notify) CreateNotify() error {
	return NotifyTable.Create(noti)
}

func UpdateRead(idNoti string) error {
	var up = bson.M{"is_read": true}
	return NotifyTable.UnsafeUpdateByID(idNoti, up)
}

func GetAllNotifyByCus(idCus string, skip int, limit int) (noties []*Notify) {
	var query = bson.M{"customer_id": idCus}
	var err = NotifyTable.Find(query).Sort("-created_at").Skip(skip).Limit(limit).All(&noties)
	rest.IsErrorRecord(err)
	rest.AssertNil(err)
	return
}
