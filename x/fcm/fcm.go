package fcm

import (
	"cetm_booking/x/mlog"
	fcm "github.com/NaySoftware/go-fcm"
)

var logFCM = mlog.NewTagLog("FCM")

const (
	RESPONSE_FAIL = "fail"
)

type FcmClient struct {
	*fcm.FcmClient
}

type FmcMessage struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type FcmMessageData struct {
	Data interface{} `json:"data,omitempty"`
	FmcMessage
}

func NewFCM(serverKey string) *FcmClient {
	return &FcmClient{
		FcmClient: fcm.NewFcmClient(serverKey),
	}
}

func (f *FcmClient) SendToMany(ids []string, data FmcMessage) (error, string) {
	var noti = fcm.NotificationPayload{
		Title: data.Title,
		Body:  data.Body,
	}
	f.NewFcmRegIdsMsg(ids, data)
	f.SetNotificationPayload(&noti)
	status, err := f.Send()
	if err != nil {
		logFCM.Debugln(0, err)
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOne(id string, data FmcMessage) (error, string) {
	return f.SendToMany([]string{id}, data)
}

func (f *FcmClient) SendToManyData(ids []string, data FcmMessageData) (error, string) {
	var noti = fcm.NotificationPayload{
		Title: data.Title,
		Body:  data.Body,
	}
	f.NewFcmRegIdsMsg(ids, data.Data)
	f.SetNotificationPayload(&noti)
	status, err := f.Send()
	if err != nil {
		logFCM.Debugln(0, err)
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOneData(id string, data FcmMessageData) (error, string) {
	return f.SendToManyData([]string{id}, data)
}
