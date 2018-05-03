package ticket

import (
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	var userTK, push = auth.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.CustomerID = userTK.ID
	var extra, _ = json.Marshal(map[string]interface{}{
		"ticket":     body,
		"push_token": push.PushToken})

	var ticket = ActionChange("", userTK.ID, ticket_onl.BOOKING_STATE_CREATED, extra)
	var countPP int
	if ticket.TypeTicket == ticket_onl.TYPE_NOW {
		countPP = CreateNumCetm(userTK, ticket)
	} else {
		var tks, _ = ticket_onl.GetAllTicketByTimeSearch(body.TimeGoBank, body.TypeTicket)
		if tks != nil {
			countPP = len(tks)
		} else {
			countPP = 0
		}
		var timeNow = math.GetTimeNowVietNam()
		if math.CompareDayTime(timeNow, body.TimeGoBank) == 0 {
			UpdateCounterTkCetm(userTK, ticket)
		}
	}
	var res = resData{
		TicketBooking: ticket,
		CountPeople:   countPP,
	}
	s.SendData(ctx, res)
}

var errOutTicket = errors.New("Đã hết chỗ trong thời gian này! Vui lòng đặt lại thời gian!")

func ValidateTicket(body ticket_onl.TicketBookingCreate) (err error) {
	var serviceID = body.ServiceID
	var btks = SetBankTickets(body.BranchID, serviceID, body.TimeGoBank, body.TimeGoBank)
	var tickets = btks.Tickets
	var countTkNow int
	var tkOrthers = make([]string, 0)
	for _, tk := range tickets {
		if serviceID == tk.ServiceID {
			countTkNow++
		} else {
			tkOrthers = append(tkOrthers, tk.ServiceID)
		}
	}
	if countTkNow >= btks.Bank.CountCounterService {
		err = errOutTicket
		return
	} else if len(tkOrthers) > 0 {
		var arrCounter = make(map[string]string, 0)
		for _, item := range btks.Bank.ServiceInCounters {
			arrCounter[item] = item
		}
		var countOSerOut int
		var countOSerIn int
		for _, ser := range tkOrthers {
			for _, cou := range btks.Bank.Counters {
				var couId = cou.ID
				for _, serC := range cou.Services {
					if ser == serC {
						if _, ok := arrCounter[couId]; ok {
							countOSerIn++
						} else {
							countOSerOut++
						}
						break
					}
				}
			}
		}
		var lenCounter = len(btks.Bank.Counters)
		if countOSerOut+countOSerIn >= lenCounter || countTkNow+countOSerOut >= lenCounter || countTkNow+countOSerIn >= btks.Bank.CountCounterService {
			err = errOutTicket
			return
		}
	}

	return
}
