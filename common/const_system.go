package common

const (
	ERROR_DONT_ROLE = "Quyền không hợp lệ !"
	NOT_EXIST       = "not found" // không có data
)

type TypePayment int

const (
	TYPE_PAYMENT_CARD  = TypePayment(1) // kiểu thanh toán theo thẻ
	TYPE_PAYMENT_MONEY = TypePayment(2) //Thanh toán theo trả lương
)
