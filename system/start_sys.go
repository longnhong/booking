package system

import (
	"cetm_booking/common"
	"cetm_booking/o/push_token"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/mlog"
	"cetm_booking/x/timer"
	"fmt"
	"time"
)

var logSys = mlog.NewTagLog("System")

func Start() (tkWorker *TicketWorker) {
	tkWorker = newCacheTicketWorker()
	tkWorker.SetCacheTicketDay()
	return
}
func (tkWorker *TicketWorker) Launch() {
	go tkWorker.startCache()
}

func CycleDayMissed() {
	ticket_onl.UpdateMissedTickets()
}

func (tkWorker *TicketWorker) SetCacheTicketDay() {
	var tickets, _ = ticket_onl.GetAllTicketDay()
	var timeNow = math.GetTimeNowVietNam()
	for _, item := range tickets {
		if math.CompareDayTime(timeNow, item.TimeGoBank) == 0 {
			var hour = math.HourMinuteEpoch(item.TimeGoBank)

			var tkDay = ticket_onl.TicketDay{
				TicketBooking: item,
				HourTimeGo:    hour,
			}
			if _, ok := tkWorker.TicketCaches[tkDay.ID]; !ok {
				tkWorker.TicketCaches[tkDay.ID] = &tkDay
			}
			continue
		}
	}
	fmt.Printf("Số ticket trong ngày", len(tkWorker.TicketCaches))
}

func (tkWorker *TicketWorker) getTicketSenPush() {
	fmt.Println("========PUSH===========")
	var tkDays = tkWorker.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		go sendPush(timeNow, tk, 0)
	}
}

func sendPush(timeNow float32, tk *ticket_onl.TicketDay, typeSend int) {
	switch typeSend {
	case 0:
		var timeRes = float64(tk.HourTimeGo - timeNow)
		if !tk.IsUsedPush && common.ConfigSystemBooking.SendNotifyBfHour >= timeRes && timeRes >= 0 &&
			tk.Status == ticket_onl.BookingStateCreated && float64(timeNow) >= common.ConfigSystemBooking.SendNotifyStartDay {
			var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
			if len(cus) > 0 {
				err := sendPushTicketDay(cus, tk.TicketBooking)
				if err == nil {
					tk.IsUsedPush = true
				}
			}
		}
	case 1:
		var timeRes = float64(tk.HourTimeGo - timeNow)
		if !tk.IsUsedNear && common.ConfigSystemBooking.StartNear >= timeRes && timeRes >= common.ConfigSystemBooking.EndNear &&
			tk.Status == ticket_onl.BookingStateCreated && tk.TypeTicket == ticket_onl.TypeSchedule {
			var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
			fmt.Printf("Số push", cus)
			if len(cus) > 0 {
				err := sendPushTicketNear(cus, tk.TicketBooking)
				if err == nil {
					tk.IsUsedNear = true
				}
			}
		}
	case 2:
		var cus, _ = push_token.GetPushsUserId(tk.CustomerID)
		fmt.Printf("Số push", cus)
		if len(cus) > 0 {
			err := sendPushTicketOut(cus, tk.TicketBooking)
			if err == nil {
				tk.IsUsedOut = true
			}
		}
		tk.Status = ticket_onl.BookingStateNotArrived
	default:
		return
	}

}

func (tkWorker *TicketWorker) getTicketSenPushNear() {
	fmt.Println("========NEAR===========")
	var tkDays = tkWorker.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		go sendPush(timeNow, tk, 1)
	}
}

func (tkWorker *TicketWorker) getTicketSenPushOut() (ticketDays []*ticket_onl.TicketDay, ids []string) {
	ticketDays = make([]*ticket_onl.TicketDay, 0)
	fmt.Println("========OUT===========")
	var tkDays = tkWorker.TicketCaches
	var timeNow = math.HourMinute()
	for _, tk := range tkDays {
		var timeRes = float64(timeNow - tk.HourTimeGo)
		if !tk.IsUsedNear && common.ConfigSystemBooking.StartOut > timeRes && timeRes >= common.ConfigSystemBooking.EndOut &&
			tk.Status == ticket_onl.BookingStateCreated && tk.TypeTicket == ticket_onl.TypeSchedule {
			ticketDays = append(ticketDays, tk)
			ids = append(ids, tk.ID)
		}
	}
	return
}

func (tkWorker *TicketWorker) sendPushOut() {
	var tkDays, ids = tkWorker.getTicketSenPushOut()
	for _, tk := range tkDays {
		go sendPush(0, tk, 2)
	}
	if len(ids) > 0 {
		ticket_onl.UpdateStatusTickets(ids, ticket_onl.BookingStateNotArrived)
	}
}

func (c *TicketWorker) startCache() {
	every15Minute := time.Tick(time.Duration(common.ConfigSystemBooking.CyclePushDay) * time.Minute)
	//every2Minute := time.Tick(time.Duration(common.ConfigSystemBooking.CyclePushTicket) * time.Minute)
	daily := timer.NewDailyTimer(common.ConfigSystemBooking.TimeSetCache, 0)
	daily.Start()
	for {
		select {
		case <-every15Minute:
			var timeNow = math.HourMinute()
			if timeNow > 5 && timeNow < 23 {
				c.getTicketSenPush()
				c.getTicketSenPushNear()
				c.sendPushOut()
			}
		case action := <-c.TicketUpdate:
			c.TicketWorking(action)
		case <-daily.C:
			fmt.Println("======== EVERYDAY ==========")
			CycleDayMissed()
			c.removeTksTicketWorkerDay()
			c.SetCacheTicketDay()
		}

	}
}
