package ctrl_to_cetm

import (
	"cetm_booking/common"
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/web"
	"errors"
)

func SearchBank(branchID string, serviceID string) (*InfoBank, error) {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_bank?branch_id=" + branchID + "&service_id=" + serviceID
	var data *DataBank
	var err error
	err = web.ResUrlClientGet(urlStr, &data)
	if err != nil {
		return nil, err
	}
	if data.Status == "error" {
		err = errors.New("Không tìm thấy Branch này!")
	}
	return data.Data, err
}

func UpdateCounterTkCetm(userTK *user.User, ticket *ticket_onl.TicketBooking) (err error) {
	var dataTicketSend = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/update_bticket"
	var data = struct {
		Data   interface{} `json:"data"`
		Status string      `json:"status"`
	}{}
	err = web.ResParamArrUrlClient(urlStr, dataTicketSend, &data)
	if err != nil {
		return
	}
	if data.Status == "error" {
		err = errors.New(data.Data.(string))
	}
	return
}

func CreateTicket(tk *ticket_onl.TicketBooking) (*InfoBank, error) {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/system_add_bkticket"
	var data *DataBank
	err := web.ResUrlClientGet(urlStr, &data)
	if err != nil {
		return nil, err
	}
	if data.Status == "error" {
		err = errors.New("Không tìm thấy Branch này!")
	}
	return data.Data, err
}

func CreateNumCetm(userTK *user.User, ticket *ticket_onl.TicketBooking, isCreate bool) (countPP int, err error) {
	var dataTicketSend = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/system_add_bkticket"
	var data = DataTicketBookNow{}
	err = web.ResParamArrUrlClient(url, dataTicketSend, &data)
	if err != nil {
		return 0, err
	}
	if isCreate {
		ticket.UpdateByCnumCetm(data.Data.Cnum, data.Data.Id)
	} else {
		ticket.UpdateTimeCheckIn(data.Data.Cnum, data.Data.Id)
	}
	return data.Data.CountPeople, err
}
