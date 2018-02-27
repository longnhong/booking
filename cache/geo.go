package cache

import (
	"cetm_booking/common"
	"cetm_booking/o/order"
	"g/geo/s2"
)

type entry struct {
	ord *order.Order
	s2.S2Entry
}

func (e *entry) updateLocation() {
	var origin = e.ord.AddressLoc
	if origin == nil {
		origin = &order.Location{}
	}
	var lat = float32(origin.Lat)
	var lng = float32(origin.Lng)
	e.S2Move(lat, lng)
}

type OrderMap struct {
	*s2.Map
	cache map[string]*entry
}

func (m *OrderMap) add(s *order.Order) *entry {
	var e = &entry{ord: s}
	m.cache[s.ID] = e
	e.updateLocation()
	m.Map.Insert(e)
	return e
}

func (m *OrderMap) remove(s *order.Order) {
	var e, ok = m.cache[s.ID]
	if !ok {
		return
	}
	delete(m.cache, s.ID)
	m.Map.Remove(e)
}

var globalMap = OrderMap{
	Map:   s2.NewMap(),
	cache: make(map[string]*entry),
}

func Refresh(s *order.Order) {
	if s == nil {
		return
	}

	var _, ok = globalMap.cache[s.ID]
	if ok {
		switch s.Status {
		case common.ORDER_STATUS_OPEN:
		case common.ORDER_STATUS_ACCEPTED:
		default:
			globalMap.remove(s)
		}
		return
	}
	globalMap.add(s)
}
func Remove(s *order.Order) {
	globalMap.remove(s)
}
