package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type resData struct {
	*ticket_onl.TicketBooking
	CountPeople int `json:"count_people"`
}

type resTime struct {
	ID         string                `json:"id"`
	TimeGoBank int64                 `json:"time_go_bank"`
	TypeTicket ticket_onl.TypeTicket `json:"type_ticket"`
	ServiceID  string                `json:"service_id"`
}

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	fmt.Println("CREATE TICKET")
	var userTK, push = auth.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	rest.AssertNil(validate(body))
	body.CustomerID = userTK.ID
	var extra, _ = json.Marshal(map[string]interface{}{
		"ticket":     body,
		"push_token": push.PushToken})

	var ticket = ActionChange("", userTK.ID, ticket_onl.BOOKING_STATE_CREATED, extra)
	var countPP int
	if ticket.TypeTicket == ticket_onl.TYPE_NOW {
		countPP, _ = ctrl.CreateNumCetm(userTK, ticket)
	} else {
		var tks, _ = ticket_onl.GetAllTicketByTimeSearch(body.TimeGoBank, body.TypeTicket)
		if tks != nil {
			countPP = len(tks)
		} else {
			countPP = 0
		}
		var timeNow = math.GetTimeNowVietNam()
		if math.CompareDayTime(timeNow, body.TimeGoBank) == 0 {
			ctrl.UpdateCounterTkCetm(userTK, ticket)
		}
	}
	var res = resData{
		TicketBooking: ticket,
		CountPeople:   countPP,
	}
	s.SendData(ctx, res)
}

var errOutTicket = errors.New("Đã hết chỗ trong thời gian này! Vui lòng đặt vào giờ khác!")

func ValidateTicket(body ticket_onl.TicketBookingCreate) error {
	var serviceID = body.ServiceID
	var btks, err1 = SetBankTickets(body.BranchID, serviceID, body.TimeGoBank, body.TimeGoBank)
	if err1 != nil {
		return err1
	}
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
		return errOutTicket
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
			return errOutTicket
		}
	}

	return nil
}

func validate(body ticket_onl.TicketBookingCreate) error {
	var serviceID = body.ServiceID
	var btks, err = SetBankTickets(body.BranchID, serviceID, body.TimeGoBank, body.TimeGoBank)
	if err != nil {
		return err
	}
	if len(btks.Tickets) >= btks.Bank.CountCounterService {
		return errOutTicket
	}
	return nil
}
