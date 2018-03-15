package rate

import (
	"cetm_booking/x/rest"
)

func (rt *Rate) CrateRate() *Rate {
	rest.AssertNil(RateTable.Create(rt))
	return rt
}
