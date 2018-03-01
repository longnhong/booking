package fcm

var (
	FcmCustomer *FcmClient
)
var FCM_SERVER_KEY_CUSTOMER string

func NewFcmApp(serverKeyCus string) {
	FCM_SERVER_KEY_CUSTOMER = serverKeyCus
	if FcmCustomer == nil {
		FcmCustomer = NewFCM(serverKeyCus)
	}
}
