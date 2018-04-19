package ticket

import (
	"cetm_booking/common"
	"cetm_booking/o/auth/user"
	"cetm_booking/o/notify"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"fmt"
)

func SendFeedback(pDevices []string, tk *ticket_onl.TicketBooking) {
	var title = "Feedback cho dịch vụ!"
	var des = "Teller " + tk.Teller + " đã phục vụ bạn. Hãy phản hồi về chất lượng dịch vụ!"
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
	}
	rest.AssertNil(noti.CreateNotify())
	var notifyTk = NotifyTicket{}
	notifyTk.Notify = &noti
	notifyTk.Ticket = tk
	var send = fcm.FcmMessageData{
		Data: notifyTk,
	}
	send.Title = title
	send.Body = des
	fcm.FcmCustomer.SendToManyData(pDevices, send)
}

func SendFee(pDevice string, tk *ticket_onl.TicketBooking) {
	var title = "Tạo vé thành công!"
	var des = "Ngân hàng đã trừ 5000VND từ tài khoản của bạn!"
	var noti = notify.Notify{
		Title:       title,
		Description: des,
		BticketID:   tk.ID,
		CustomerId:  tk.CustomerID,
	}
	rest.AssertNil(noti.CreateNotify())
	var notifyTk = NotifyTicket{}
	notifyTk.Notify = &noti
	notifyTk.Ticket = tk
	var send = fcm.FcmMessageData{
		Data: notifyTk,
	}
	send.Title = title
	send.Body = des
	fcm.FcmCustomer.SendToOneData(pDevice, send)
}

func CreateNumCetm(userTK *user.User, ticket *ticket_onl.TicketBooking) (countPP int) {
	var dataTicketSend = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/system_add_bkticket"
	var data = DataTicketBookNow{}
	rest.AssertNil(web.ResParamArrUrlClient(url, dataTicketSend, &data))
	ticket.UpdateByCnumCetm(data.Data.Cnum, data.Data.Id)
	fmt.Printf("Data tao ve ", data.Data)
	return data.Data.CountPeople
}
