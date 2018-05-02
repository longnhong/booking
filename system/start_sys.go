package system

import (
	"cetm_booking/o/notify"
	"cetm_booking/o/push_token"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"cetm_booking/x/mlog"
	"fmt"
	"time"
)

var logSys = mlog.NewTagLog("System")

func Launch() {
	SetCacheTicketDay()
	go startCache(TicketWorkerDay)
}

func SetCacheTicketDay() {
	var tickets, _ = ticket_onl.GetAllTicketDay()
	TicketWorkerDay = newCacheTicketWorker()
	var timeNow = math.GetTimeNowVietNam()
	for _, item := range tickets {
		if math.CompareDayTime(timeNow, item.TimeGoBank) == 0 {
			var hour = math.HourMinuteEpoch(item.TimeGoBank)

			var tkDay = ticket_onl.TicketDay{
				TicketBooking: item,
				HourTimeGo:    hour,
			}
			if _, ok := TicketWorkerDay.TicketCaches[tkDay.ID]; !ok {
				TicketWorkerDay.TicketCaches[tkDay.ID] = &tkDay
			}
			continue
		}
	}
	fmt.Printf("Số ticket trong ngày", len(TicketWorkerDay.TicketCaches))
}

func getTicketSenPush() {
	fmt.Println("========PUSH===========")
	var tkDays = TicketWorkerDay.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		var timeRes = tk.HourTimeGo - timeNow
		if !tk.IsUsedPush && 3 >= timeRes && timeRes >= 0 && tk.Status == ticket_onl.BOOKING_STATE_CREATED && timeNow >= 7 {
			var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
			if len(cus) > 0 {
				err := sendPushTicketDay(cus, tk.TicketBooking)
				if err == nil {
					tk.IsUsedPush = true
				}
			}
		}
	}
}

func getTicketSenPushNear() {
	fmt.Println("========NEAR===========")
	var tkDays = TicketWorkerDay.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		var timeRes = tk.HourTimeGo - timeNow
		if !tk.IsUsedNear && 0.5 >= timeRes && timeRes >= 0.25 && tk.Status == ticket_onl.BOOKING_STATE_CREATED && tk.TypeTicket == ticket_onl.TYPE_SCHEDUCE {
			var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
			fmt.Printf("Số push", cus)
			if len(cus) > 0 {
				err := sendPushTicketNear(cus, tk.TicketBooking)
				if err == nil {
					tk.IsUsedNear = true
				}
			}
		}
	}
}

func getTicketSenPushOut() (ticketDays []*ticket_onl.TicketDay) {
	ticketDays = make([]*ticket_onl.TicketDay, 0)
	fmt.Println("========OUT===========")
	var tkDays = TicketWorkerDay.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		var timeRes = tk.HourTimeGo - timeNow
		if !tk.IsUsedNear && -0.25 > timeRes && timeRes >= -0.5 &&
			tk.Status == ticket_onl.BOOKING_STATE_CREATED && tk.TypeTicket == ticket_onl.TYPE_SCHEDUCE {
			ticketDays = append(ticketDays, tk)
		}
	}
	return
}

func sendPushOut() {
	var tkDays = getTicketSenPushOut()
	var ids = make([]string, 0)
	for _, tk := range tkDays {
		var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
		fmt.Printf("Số push", cus)
		if len(cus) > 0 {
			err := sendPushTicketOut(cus, tk.TicketBooking)
			if err == nil {
				tk.IsUsedOut = true
			}
		}
		ids = append(ids, tk.ID)
		tk.Status = ticket_onl.BOOKING_STATE_NOT_ARRIVED
	}
	if len(ids) > 0 {
		ticket_onl.UpdateStatusTickets(ids, ticket_onl.BOOKING_STATE_NOT_ARRIVED)
	}
}

func startCache(c *ticketWorker) {
	every15Minute := time.Tick(1 * time.Minute)
	every2Minute := time.Tick(2 * time.Minute)
	for {
		select {
		case <-every15Minute:
			getTicketSenPush()
		case <-every2Minute:
			getTicketSenPushNear()
			sendPushOut()
		case action := <-c.TicketUpdate:
			c.TicketWorking(action)
		}
	}
}

func sendPushTicketDay(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Bạn có lịch hẹn ngày hôm nay!"
	var des = "Kính chào quý khách! Hôm nay bạn có lịch hẹn." +
		"\nThời gian: " + math.TimeToString(tk.TimeGoBank) +
		"\nĐịa chỉ: " + tk.BranchAddress +
		"\nQuý khách vui lòng đến trước hoặc sau 15p để được ưu tiên phục vụ! Xin cảm ơn."
	err = setDataNoti(title, des, tk, pDevice, notify.CETM_SCHEDULE_DAY)
	return
}

func sendPushTicketNear(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Sắp đến giờ phục vụ!"
	var des = "Kính chào quý khách! Thời gian " + math.TimeToString(tk.TimeGoBank-15*60) +
		" đến " + math.TimeToString(tk.TimeGoBank+15*60) + " quý khách sẽ được ưu tiên phục vụ. Hãy lại quầy giao dịch để thực hiện yêu cầu của bạn.\nXin cảm ơn!"

	err = setDataNoti(title, des, tk, pDevice, notify.CETM_NEAR_HOUR)
	return
}

func sendPushTicketOut(pDevice []string, tk *ticket_onl.TicketBooking) (err error) {
	var title = "Lịch hẹn đã bị hủy!"
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
