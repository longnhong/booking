package system

import (
	"cetm_booking/o/ticket_onl"
	"encoding/json"
)

func (action *TicketAction) actionCheckCode() {
	var data = struct {
		CustomerCode string `json:"customer_code"`
		BranchID     string `json:"branch_id"`
	}{}
	var err1 = json.Unmarshal(action.Extra, &data)
	if err1 != nil {
		action.SetError(err1)
		return
	}
	var ticket, err = ticket_onl.CheckCustomerCode(data.CustomerCode, data.BranchID)
	if err != nil {
		action.SetError(err)
		return
	}
	action.Ticket = ticket
}
