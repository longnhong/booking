package system

import (
	"cetm_booking/common"
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/ultility"
	"encoding/json"
	"errors"
)

func (action *TicketAction) actionCreateCetm() {
	var btks, err = ticket_onl.GetTicketDayByUser(action.CusID)
	if err != nil {
		action.SetError(err)
		return
	}
	var timeNow = math.HourMinute()
	var ticket *ticket_onl.TicketBooking
	if len(btks) > 0 {
		for _, item := range btks {
			var timeGoBank = math.HourMinuteEpoch(item.TimeGoBank)
			var resTime = float64(timeNow - timeGoBank)
			if resTime == 0 || (resTime > 0 && resTime < common.ConfigSystemBooking.EndNear) ||
				(resTime < 0 && resTime > -common.ConfigSystemBooking.EndNear) {
				ticket = item
			}
		}
	}
	if ticket == nil {
		err = errors.New("Vé của bạn chưa đến hoặc quá giờ! Vui lòng kiếm tra lại!")
		action.SetError(err)
		return
	}
	var data *common.Location
	err = json.Unmarshal(action.Extra, &data)
	if err != nil {
		action.SetError(err)
		return
	}
	bank, err := ctrl.SearchBank(ticket.BranchID, ticket.ServiceID)
	if err != nil {
		action.SetError(err)
		return
	}
	var scanKm = ultility.Haversine(data.Lat, data.Lng, bank.Lat, bank.Lng)
	if scanKm > common.ConfigSystemBooking.ScanNear {
		err = errors.New("Bạn đang không trong phạm vi ngân hàng!")
		action.SetError(err)
		return
	}
	action.Ticket = ticket
}
