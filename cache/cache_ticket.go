package cache

import (
	"cetm_booking/o/push_token"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"strconv"
	"time"
)

var CacheTicketAll = make(map[int][]*ticket_onl.TicketBooking, 0)

func AddCachePushCustomer() {
	var timeNow = time.Now()
	var tickets, _ = ticket_onl.GetAllTicketDay()
	var mapTk = make(map[int][]*ticket_onl.TicketBooking, 0)

	for _, item := range tickets {
		var tkAdd = make([]*ticket_onl.TicketBooking, 0)
		if math.CompareDayTime(timeNow, item.CheckInAt) == 0 {
			var hourItem = 7
			var dateTicket = time.Unix(item.CheckInAt, 0)
			if dateTicket.Hour() > 10 {
				hourItem = dateTicket.Hour() - 3
			}
			if tks, ok := mapTk[hourItem]; ok {
				tks = append(tks, item)
				mapTk[hourItem] = tks
			} else {
				tkAdd = append(tkAdd, item)
				mapTk[int(hourItem)] = tkAdd
			}
			break
		}
	}
	CacheTicketAll = mapTk
}

func GetTicketDayAndSendPush() {
	AddCachePushCustomer()
	go SendPush()
}

func SendPush() {
	for {
		var timeHour = time.Now().Hour()
		if len(CacheTicketAll) > 0 {
			if tks, ok := CacheTicketAll[timeHour]; ok {
				for _, tk := range tks {
					var dateTicket = time.Unix(tk.CheckInAt, 0)
					var pushTokens, _ = push_token.GetPushsUserId(tk.CustomerID)
					var noti = fcm.FmcMessage{
						Title: "Gần đến giờ đặt vé!",
						Body:  "Thời gian: " + strconv.Itoa(dateTicket.Hour()) + ":" + strconv.Itoa(dateTicket.Minute()) + ". \n Trân trọng kính mời quý khách!",
					}
					fcm.FcmCustomer.SendToMany(pushTokens, noti)
				}
			}
		}
	}
}
