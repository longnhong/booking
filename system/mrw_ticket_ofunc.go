package system

import (
	"cetm_booking/o/notify"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/rest"
)

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
	rest.AssertNil(noti.CreateNotify())
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
	rest.AssertNil(noti.CreateNotify())
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
