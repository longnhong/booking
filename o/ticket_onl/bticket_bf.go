package ticket_onl

import (
	"cetm_booking/x/rest"
	"cetm_booking/x/rest/math"
	"errors"
)

var errBranchID = rest.BadRequestValid(errors.New("Không tồn tại chi nhánh"))
var errService = rest.BadRequestValid(errors.New("Không tồn tại service"))
var errCustomer = rest.BadRequestValid(errors.New("Không lấy được thông tin khách hàng"))
var errTimeGoBank = rest.BadRequestValid(errors.New("Chọn thời gian đến ngân hàng"))
var errBTicketID = rest.BadRequestValid(errors.New("Chọn vé trước khi sửa thông tin"))

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
		ServiceID:    btbk.ServiceID,
		TimeGoBank:   btbk.TimeGoBank,
		CustomerCode: cusCode,
		Lang:         lang,
		Status:       BookingStateCreated,
	}
	return nil, &ticket
}

func (btbk *TicketUpdate) updateBf() error {
	var err = btbk.CheckTicketBooking()
	if err != nil {
		return err
	}
	return nil
}

func (tc *TicketUpdate) CheckTicketBookingUp() error {
	if len(tc.BTicketID) == 0 {
		return errBTicketID
	}
	err := tc.TicketBookingCreate.CheckTicketBooking()
	if err != nil {
		return err
	}
	return nil
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
	return nil
}
