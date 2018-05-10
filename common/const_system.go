package common

const (
	ERROR_DONT_ROLE = "Quyền không hợp lệ !"
	NOT_EXIST       = "not found" // không có data
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
