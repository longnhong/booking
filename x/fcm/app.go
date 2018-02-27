package fcm

var (
	FcmCustomer *FcmClient
	FcmEmployee *FcmClient
)
var FCM_SERVER_KEY_CUSTOMER string
var FCM_SERVER_KEY_EMPLOYEE string

func NewFcmApp(serverKeyCus string, serverKeyEmp string) {
	FCM_SERVER_KEY_CUSTOMER = serverKeyCus
	if FcmCustomer == nil {
		FcmCustomer = NewFCM(serverKeyCus)
	}
	FCM_SERVER_KEY_EMPLOYEE = serverKeyEmp
	if FcmEmployee == nil {
		FcmEmployee = NewFCM(serverKeyEmp)
	}
}
