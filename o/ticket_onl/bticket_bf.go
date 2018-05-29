package ticket_onl

import (
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"errors"
)

var errBranchID = rest.BadRequestValid(errors.New("Không tồn tại chi nhánh"))
var errService = rest.BadRequestValid(errors.New("Không tồn tại service"))
var errCustomer = rest.BadRequestValid(errors.New("Không lấy được thông tin khách hàng"))
var errTimeGoBank = rest.BadRequestValid(errors.New("Chọn thời gian đến ngân hàng"))
var errBTicketID = rest.BadRequestValid(errors.New("Chọn vé trước khi sửa thông tin"))
var errTypeTicket = rest.BadRequestValid(errors.New("Không có loại ticket này!"))

func NewParamDefault() (customerCode string) {
	customerCode = math.RandNumString(6)
	return
}

func (btbk *TicketBookingCreate) createBf() (error, *TicketBooking) {
	var err = btbk.CheckTicketBooking()
	if err != nil {
		return err, nil
	}

	var cusCode = NewParamDefault()
	var ticket = TicketBooking{
		BranchID:      btbk.BranchID,
		BranchName:    btbk.BranchName,
		BranchAddress: btbk.BranchAddress,
		CustomerID:    btbk.CustomerID,
		ServiceID:     btbk.ServiceID,
		ServiceName:   btbk.ServiceName,
		TimeGoBank:    btbk.TimeGoBank,
		TypeTicket:    btbk.TypeTicket,
		CustomerCode:  cusCode,
		Lang:          btbk.Lang,
		Status:        BookingStateCreated,
	}

	if btbk.TypeTicket == TypeNow {
		ticket.CheckInAt = btbk.TimeGoBank
	}
	return nil, &ticket
}

func CheckType(typeTK TypeTicket) error {
	if typeTK != TypeNow && typeTK != TypeSchedule {
		return errTypeTicket
	}
	return nil
}

func (tk *TicketBooking) updateBf(btbk *TicketUpdate) error {
	var timeNow = math.GetTimeNowVietNam().Unix()
	btbk.Tracks = tk.Tracks
	btbk.Tracks = tk.updateTrack(btbk.ServiceID, btbk.BranchID, btbk.Status, timeNow)
	var err = btbk.CheckTicketBookingUp()
	if err != nil {
		return err
	}
	btbk.UpdatedAt = timeNow
	return nil
}

func (tk *TicketBooking) updateTrack(serviceID string, branchID string, status BookingState, timeNow int64) []TicketHst {
	var tracks = tk.Tracks
	var hst = TicketHst{
		BranchID:  branchID,
		ServiceID: serviceID,
		Status:    status,
		MTime:     timeNow,
	}
	tracks = append(tracks, hst)
	return tracks
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
	// if len(tc.Customer) == 0 {
	// 	return errCustomer
	// }
	if len(tc.ServiceID) == 0 {
		return errService
	}
	if tc.TimeGoBank == 0 {
		return errTimeGoBank
	}
	return CheckType(tc.TypeTicket)
}
