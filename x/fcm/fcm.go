package fcm

import (
	fcm "github.com/NaySoftware/go-fcm"
)

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
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOne(id string, data FmcMessage) (error, string) {
	return f.SendToMany([]string{id}, data)
}
