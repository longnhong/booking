package ticket_onl

import (
	"cetm_booking/x/rest"
	"cetm_booking/x/rest/math"
	"errors"
	"time"
)

var errBranchID = rest.BadRequestValid(errors.New("Không tồn tại chi nhánh"))
var errService = rest.BadRequestValid(errors.New("Không tồn tại service"))
var errCustomer = rest.BadRequestValid(errors.New("Không lấy được thông tin khách hàng"))
var errTimeGoBank = rest.BadRequestValid(errors.New("Chọn thời gian đến ngân hàng"))
var errBTicketID = rest.BadRequestValid(errors.New("Chọn vé trước khi sửa thông tin"))
var errTypeTicket = rest.BadRequestValid(errors.New("Không có loại ticket này!"))

func NewParamDefault() (customerCode string, lang string) {
	customerCode = math.RandNumString(6)
	lang = "vi"
	return
}

func (btbk *TicketBookingCreate) createBf() (error, *TicketBooking) {
	var err = btbk.CheckTicketBooking()
	if err != nil {
		return err, nil
	}
	var cusCode, lang = NewParamDefault()
	var ticket = TicketBooking{
		BranchID:     btbk.BranchID,
		Customer:     btbk.Customer,
		CustomerID:   btbk.CustomerID,
		ServiceID:    btbk.ServiceID,
		TimeGoBank:   btbk.TimeGoBank,
		TypeTicket:   btbk.TypeTicket,
		CustomerCode: cusCode,
		Lang:         lang,
		Status:       BOOKING_STATE_CREATED,
	}
	if btbk.TypeTicket == TYPE_NOW {
		ticket.CheckInAt = btbk.TimeGoBank
	}
	return nil, &ticket
}

func CheckType(typeTK TypeTicket) error {
	if typeTK != TYPE_NOW && typeTK != TYPE_SCHEDUCE {
		return errTypeTicket
	}
	return nil
}

func (btbk *TicketUpdate) updateBf() error {
	var err = btbk.CheckTicketBookingUp()
	if err != nil {
		return err
	}
	btbk.UpdatedAt = time.Now().Unix()
	return nil
}

func (tc *TicketUpdate) CheckTicketBookingUp() error {
	if len(tc.BTicketID) == 0 {
		return errBTicketID
	}
	if len(tc.BranchID) == 0 {
		return errBranchID
	}
	if len(tc.ServiceID) == 0 {
		return errService
	}
	if tc.TimeGoBank == 0 {
		return errTimeGoBank
	}
	return CheckType(tc.TypeTicket)
}

func (tc *TicketBookingCreate) CheckTicketBooking() error {
	if len(tc.BranchID) == 0 {
		return errBranchID
	}
	if len(tc.Customer) == 0 {
		return errCustomer
	}
	if len(tc.ServiceID) == 0 {
		return errService
	}
	if tc.TimeGoBank == 0 {
		return errTimeGoBank
	}
	return CheckType(tc.TypeTicket)
}
