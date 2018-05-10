package system

import (
	"cetm_booking/o/notify"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"cetm_booking/x/utils"
	"fmt"
)

func (c *ticketWorker) TicketWorking(action *TicketAction) error {
	fmt.Printf("\nTICKET ACTION", action)
	if action == nil {
		return nil
	}
	defer utils.Recover()
	defer action.Done()
	action.handlerAction()
	var err = action.GetError()
	if err == nil {
		var timegoBank = action.Ticket.TimeGoBank
		if action.Action == ticket_onl.BOOKING_STATE_CREATED && math.CompareDayTime(math.GetTimeNowVietNam(), timegoBank) == 0 {
			var hourDay = math.HourMinuteEpoch(timegoBank)
			fmt.Printf("\nTẠO VÉ TRONG NGÀY", hourDay)
			var tkDay = ticket_onl.TicketDay{
				TicketBooking: action.Ticket,
				HourTimeGo:    hourDay,
			}
			c.TicketCaches[action.Ticket.ID] = &tkDay
		}
	}
	return err
}

func sendFeedback(pDevices []string, tk *ticket_onl.TicketBooking, status ticket_onl.BookingState) {
	var title = "Feedback cho dịch vụ!"
	var des string
	var stateNotify notify.StateNotify
	if status == ticket_onl.BOOKING_STATE_FINISHED {
		des = "Teller " + tk.Teller + " đã phục vụ bạn. Hãy phản hồi về chất lượng dịch vụ!"
		stateNotify = notify.CETM_FINISHED
	} else {
		des = "Vé của bạn đã bị hủy. Bạn có muốn phản hồi về chất lượng dịch vụ?"
		stateNotify = notify.CETM_CANCELLED
	}
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
		State:       stateNotify,
	}
	noti.CreateNotify()
	var notifyTk = ticket_onl.NotifyTicket{}
	notifyTk.Notify = &noti
	notifyTk.Ticket = tk
	var send = fcm.FcmMessageData{
		Data: notifyTk,
	}
	send.Title = title
	send.Body = des
	fcm.FcmCustomer.SendToManyData(pDevices, send)
}

func sendFee(pDevice string, tk *ticket_onl.TicketBooking) {
	var title = "Tạo vé thành công!"
	var des = "Ngân hàng đã trừ 5000VND từ tài khoản của bạn!"
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
		State:       notify.CETM_CREATE,
	}
	noti.CreateNotify()
	var notifyTk = ticket_onl.NotifyTicket{}
	notifyTk.Notify = &noti
	notifyTk.Ticket = tk
	var send = fcm.FcmMessageData{
		Data: notifyTk,
	}
	send.Title = title
	send.Body = des
	fcm.FcmCustomer.SendToOneData(pDevice, send)
}
