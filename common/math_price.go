package common

import (
	"fmt"
	//"sort"
	"cetm_booking/o/tool"
	"cetm_booking/x/rest"
	"cetm_booking/x/rest/math"
)

type MathPriceOrder struct {
	TypeWork     TypeWork     `bson:"type_work" json:"type_work" validate:"required"`
	Promotions   []string     `bson:"promotions" json:"promotions"`
	ServiceWorks []string     `bson:"service_works" json:"service_works" validate:"required"`
	ToolServices []string     `bson:"tool_services" json:"tool_services"`
	DayWeeks     DayWeeks     `bson:"day_weeks" json:"day_weeks"`
	PeopleOther  *PeopleOther `bson:"people_other" json:"people_other"`
	DayStartWork int64        `bson:"day_start_work" json:"day_start_work" validate:"required"`
}

type PeopleOther struct {
	Phone string `bson:"phone" json:"phone"`
	Name  string `bson:"name" json:"name"`
}

type DayWeek struct {
	IdItem    string          `bson:"id_item" json:"id_item" validate:"required"`
	DateIn    int64           `bson:"date_in" json:"date_in" validate:"required"` // 2,,3,4,5,6,7,8
	HourStart float32         `bson:"hour_start" json:"hour_start" validate:"required"`
	HourEnd   float32         `bson:"hour_end" json:"hour_end" validate:"required"`
	HourDay   float32         `bson:"hour_day" json:"hour_day" validate:"required"`
	MTime     int64           `bson:"mtime" json:"mtime"`
	Status    ItemOrderStatus `bson:"status" json:"status" validate:"required"`
}

type DayWeeks []*DayWeek

func (a DayWeeks) Len() int           { return len(a) }
func (a DayWeeks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DayWeeks) Less(i, j int) bool { return a[i].DateIn < a[j].DateIn }

func (ord *MathPriceOrder) MathHourWork() float32 {
	// sort.Sort(DayWeeks(ord.DayWeeks))

	var hourInWeek float32
	for _, item := range ord.DayWeeks {
		item.IdItem = math.RandStringUpper("", 6)
		var hourDay = item.HourEnd - item.HourStart
		item.HourDay = hourDay
		item.Status = ITEM_ORDER_STATUS_NEW
		hourInWeek += hourDay
	}
	fmt.Println(hourInWeek)
	return hourInWeek
}

func (ord *MathPriceOrder) MathPriceOrder(hourMoney int) (hourInWeek float32, priceAllHour float32, priceWeek float32, priceTool float32, priceMoneyMonth float32) {
	hourInWeek = ord.MathHourWork()
	priceAllHour = hourInWeek * float32(hourMoney)
	if len(ord.ToolServices) > 0 {
		var tools, err = tool.GetToolByArrayID(ord.ToolServices)
		if err != nil && err.Error() != NOT_EXIST {
			rest.AssertNil(rest.WrapBadRequest(err, ""))
		}
		for _, item := range tools {
			priceTool += float32(item.Price)
		}
	}
	priceWeek = priceAllHour + priceTool
	priceMoneyMonth = priceWeek * 4
	return
}
