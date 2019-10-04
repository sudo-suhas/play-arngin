package gen

import (
	"reflect"

	"github.com/gofrs/uuid"
	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/dataset"
)

// Rules generates and returns n rules.
func Rules(n int) []arngin.AddonRule {
	rand := newRandHelper()
	rules := make([]arngin.AddonRule, 0, n)
	for i := 0; i < n; i++ {
		var r arngin.AddonRule
		if rand.Bool() {
			r.Sources = rand.SampleInts(dataset.AddonCities, 15)
		}

		if rand.Bool() {
			r.Destinations = rand.SampleInts(dataset.AddonCities, 15)
		}

		if rand.Bool() && rand.Bool() {
			r.BoardingPts = busPtsForCities(rand, r.Sources)
		}

		if rand.Bool() && rand.Bool() {
			r.DroppingPts = busPtsForCities(rand, r.Destinations)
		}

		if rand.Bool() && rand.Bool() {
			start := rand.Floatn(20)
			r.BoardingTime = &arngin.FloatRange{GTE: start, LTE: start + 3.5}
		}

		if rand.Bool() && rand.Bool() {
			start := rand.Floatn(20)
			r.DroppingTime = &arngin.FloatRange{GTE: start, LTE: start + 3.5}
		}

		if rand.Bool() && rand.Bool() {
			r.SeatCount = randNumRule(rand, 10)
		}

		if rand.Bool() && rand.Bool() {
			if rand.Bool() && rand.Bool() && rand.Bool() {
				r.BusOperators = rand.SampleStrings(dataset.BusOperators, 10)
			} else {
				r.BusOperators = rand.SampleStrings(dataset.PrimeBusOperators, 10)
			}
		}

		if rand.Bool() && rand.Bool() {
			r.Duration = randNumRule(rand, 24)
		}

		if rand.Bool() && rand.Bool() {
			r.Channels = rand.SampleStrings(dataset.Channels, len(dataset.Channels))
		}
		for _, c := range r.Channels {
			if c != "MOBILE_APP" {
				continue
			}
			if rand.Bool() && rand.Bool() {
				r.Appversion = randNumRule(rand, 65000)
			}
			break
		}

		if reflect.ValueOf(r).IsZero() {
			i--
			continue
		}

		r.ID = uuid.Must(uuid.NewV4())
		r.Addons = rand.SampleStrings(dataset.Addons, 10)

		rules = append(rules, r)
	}
	return rules
}

func busPtsForCities(rand randHelper, cities []int) []int {
	var pts []int
	for _, c := range cities {
		pts = append(pts, dataset.CityBusPts[c]...)
	}
	return rand.SampleInts(pts, 40)
}

func randNumRule(rand randHelper, n int) *arngin.NumRule {
	var op arngin.CompareOp
	switch rand.Intn(3) {
	case 0:
		op = arngin.EqOp
	case 1:
		op = arngin.GTEOp
	case 2:
		op = arngin.LTEOp
	}
	v := rand.Intn(n)
	return &arngin.NumRule{Op: op, Value: v}
}
