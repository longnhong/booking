package system

import (
	"cetm_booking/o/push_token"
	"cetm_booking/o/ticket_onl"
	"encoding/json"
	"errors"
	"fmt"
)

func (action *TicketAction) handlerAction() {
	var ticket *ticket_onl.TicketBooking
	var err error
	if action.Action != ticket_onl.BOOKING_STATE_CREATED && action.Action != ticket_onl.BOOKING_STATE_CHECK_CODE {
		ticket, err = GetTicketByID(action.TicketID)
		if err != nil {
			action.SetError(err)
			return
		}
	}
	switch action.Action {
	case ticket_onl.BOOKING_STATE_CREATED:
		var data = struct {
			PushToken string                         `json:"push_token"`
			Ticket    ticket_onl.TicketBookingCreate `json:"ticket"`
		}{}
		err = json.Unmarshal(action.Extra, &data)
		if err != nil {
			action.SetError(err)
			return
		}
		fmt.Printf("TICKET CREATE", data)
		ticket, err := data.Ticket.CrateTicketBooking()
		if err != nil {
			action.SetError(err)
			return
		}
		action.Ticket = ticket
		go sendFee(data.PushToken, ticket)
	case ticket_onl.BOOKING_STATE_CONFIRMED:
	case ticket_onl.BOOKING_STATE_DELETE:
		err = ticket.MarkDeleteTicket()
		if err != nil {
			action.SetError(err)
			return
		}
		ticket.Status = ticket_onl.BOOKING_STATE_DELETE
		action.Ticket = ticket
	case ticket_onl.BOOKING_STATE_CANCELLED:
		fallthrough
	case ticket_onl.BOOKING_STATE_FINISHED:
		var data *ticket_onl.UpdateCetm
		var err = json.Unmarshal(action.Extra, &data)
		if err != nil {
			action.SetError(err)
			return
		}
		if ticket.Status == ticket_onl.BOOKING_STATE_CANCELLED || ticket_onl.BOOKING_STATE_FINISHED == ticket.Status {
			action.SetError(errors.New("Vé đã được phản hồi"))
			return
		}
		err = data.UpdateTicketBookingByCetm()
		if err != nil {
			action.SetError(err)
			return
		}
		ticket.Status = data.Status
		ticket.AvatarTeller = data.AvatarTeller
		ticket.TellerID = data.TellerID
		ticket.IdTicketCetm = data.IdTicketCetm
		ticket.CnumCetm = data.CnumCetm
		ticket.ServingTime = data.ServingTime
		ticket.WaitingTime = data.WaitingTime
		ticket.Teller = data.Teller
		action.Ticket = ticket

		pDevices, err := push_token.GetPushsUserId(ticket.CustomerID)
		if err != nil {
			action.SetError(err)
			return
		}
		go sendFeedback(pDevices, ticket, action.Action)
	case ticket_onl.BOOKING_STATE_NOT_ARRIVED:
	case ticket_onl.BOOKING_STATE_CHECK_CODE:
		var data = struct {
			CustomerCode string `json:"customer_code"`
			BranchId     string `json:"branch_id"`
		}{}
		var err1 = json.Unmarshal(action.Extra, &data)
		if err1 != nil {
			action.SetError(err1)
			return
		}
		var ticket, err = ticket_onl.CheckCustomerCode(data.CustomerCode, data.BranchId)
		if err != nil {
			action.SetError(err)
			return
		}
		ticket.UpdateTimeCheckIn()
		action.Ticket = ticket
		if val, ok := TicketWorkerDay.TicketCaches[ticket.ID]; ok {
			val.TicketBooking = ticket
		}
	default:
		err := errors.New("No Action")
		action.SetError(err)
	}
}
