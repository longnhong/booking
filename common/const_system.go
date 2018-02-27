package common

const (
	ERROR_DONT_ROLE = "Quyền không hợp lệ !"
	NOT_EXIST       = "not found" // không có data
)

var 

type OrderStatus string
type ItemOrderStatus string

const (
	ORDER_STATUS_OPEN          = OrderStatus("open")
	ORDER_STATUS_BIDDING       = OrderStatus("bidding")
	ORDER_STATUS_ACCEPTED      = OrderStatus("accepted")
	ORDER_STATUS_WORKING       = OrderStatus("working")
	ORDER_STATUS_FINISHED      = OrderStatus("finished")
	ORDER_STATUS_CANCELED      = OrderStatus("canceled")
	ITEM_ORDER_STATUS_NEW      = ItemOrderStatus("new")
	ITEM_ORDER_STATUS_WORKING  = ItemOrderStatus("working")
	ITEM_ORDER_STATUS_FINISHED = ItemOrderStatus("finished")
)

type TypeWork int

const (
	TYPE_DAY   = TypeWork(1) // kiểu làm việc theo ngày
	TYPE_WEEK  = TypeWork(2) // kiểu làm việc theo tuần
	TYPE_MONTH = TypeWork(3) // kiểu làm việc theo tháng
	TYPE_YEAR  = TypeWork(4) // kiểu làm việc theo năm
)

type TypePayment int

const (
	TYPE_PAYMENT_CARD  = TypePayment(1) // kiểu thanh toán theo thẻ
	TYPE_PAYMENT_MONEY = TypePayment(2) //Thanh toán theo trả lương
)
