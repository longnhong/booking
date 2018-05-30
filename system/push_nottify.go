package system

import (
	"cetm_booking/o/notify"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
)

func sendPushTicketDay(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Bạn có lịch hẹn ngày hôm nay"
	var des = "Kính chào quý khách! Hôm nay bạn có lịch hẹn." +
		"\nThời gian: " + math.TimeToString(tk.TimeGoBank) +
		"\nĐịa chỉ: " + tk.BranchAddress +
		"\nQuý khách vui lòng đến trước hoặc sau 15p để được ưu tiên phục vụ! Xin cảm ơn!"
	err = setDataNoti(title, des, tk, pDevice, notify.CETM_SCHEDULE_DAY)
	return
}

func sendPushTicketNear(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Sắp đến giờ phục vụ"
	var des = "Kính chào quý khách! Thời gian " + math.TimeToString(tk.TimeGoBank-15*60) +
		" đến " + math.TimeToString(tk.TimeGoBank+15*60) + " quý khách sẽ được ưu tiên phục vụ. Hãy lại quầy giao dịch để thực hiện yêu cầu của bạn.\nXin cảm ơn!"

	err = setDataNoti(title, des, tk, pDevice, notify.CETM_NEAR_HOUR)
	return
}

func sendPushTicketOut(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Lịch hẹn đã bị hủy"
	var des = "Kính chào quý khách! Do không có sự xác nhận làm dịch vụ theo yêu cầu của bạn. Chúng tôi đã buộc hủy lịch hẹn của bạn." +
		"\n Lên đơn mới để được phục vụ lại.\nXin cảm ơn!"
	err = setDataNoti(title, des, tk, pDevice, notify.CETM_NOT_ARRIVED)
	return
}

func setDataNoti(title string, des string, tk *ticket_onl.TicketBooking, pDevice []string, state notify.StateNotify) (err error) {
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
		State:       state,
	}
	err = noti.CreateNotify()
	if err != nil {
		return
	}
	var notifyTk = ticket_onl.NotifyTicket{}
	notifyTk.Notify = &noti
	notifyTk.Ticket = tk
	var send = fcm.FcmMessageData{
		Data: notifyTk,
	}
	send.Title = title
	send.Body = des
	err, _ = fcm.FcmCustomer.SendToManyData(pDevice, send)
	if err != nil {
		noti.RemoveNotify()
	}
	return
}

func sendFeedback(pDevices []string, tk *ticket_onl.TicketBooking, status ticket_onl.BookingState) {
	var title = "Feedback cho dịch vụ."
	var des string
	var stateNotify notify.StateNotify
	if status == ticket_onl.BookingStateFinished {
		des = "Teller " + tk.Teller + " đã phục vụ bạn. Hãy phản hồi về chất lượng dịch vụ."
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
	var title = "Tạo vé thành công"
	var des = "Ngân hàng đã trừ 5000VND từ tài khoản của bạn."
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
		State:       notify.CETM_CREATE,
	}
	noti.CreatedAt = math.GetTimeNowVietNam().Unix()
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
