package gen

import (
	arngin "github.com/sudo-suhas/play-arngin"
	"github.com/sudo-suhas/play-arngin/dataset"
)

// Qs generates and returns n addon queries.
func Qs(n int) []arngin.AddonsQ {
	rand := newRandHelper()
	var void struct{}
	qset := make(map[arngin.AddonsQ]struct{}, n)
	for i := 0; i < n; i++ {
		var q arngin.AddonsQ

		var src int
		// Only pick a low traffic city 1 in 4 times.
		if rand.Bool() && rand.Bool() {
			src = rand.IntEl(dataset.Cities)
		} else {
			src = rand.IntEl(dataset.PrimeCities)
		}
		q.Source = src

		dest := src
		if rand.Bool() && rand.Bool() {
			for dest == src {
				dest = rand.IntEl(dataset.Cities)
			}
		} else {
			for dest == src {
				dest = rand.IntEl(dataset.PrimeCities)
			}
		}
		q.Destination = dest

		// q.BoardingPt = rand.IntEl(dataset.CityBusPts[src])
		// q.DroppingPt = rand.IntEl(dataset.CityBusPts[dest])

		q.BoardingTime = rand.Floatn(23.9)
		q.DroppingTime = rand.Floatn(23.9)

		q.Seats = 1
		if rand.Bool() && rand.Bool() {
			q.Seats = rand.Intn(9) + 2
		}

		if rand.Bool() && rand.Bool() {
			q.BusOperator = rand.StringEl(dataset.BusOperators)
		} else {
			q.BusOperator = rand.StringEl(dataset.PrimeBusOperators)
		}

		q.Duration = rand.Intn(29) + 2
		q.Channel = rand.StringEl(dataset.Channels)
		if q.Channel == "MOBILE_APP" {
			q.Appversion = rand.Intn(25000) + 60000
		}

		if _, ok := qset[q]; ok {
			i--
			continue
		}
		qset[q] = void
	}

	qq := make([]arngin.AddonsQ, 0, n)
	for q := range qset {
		qq = append(qq, q)
	}

	return qq
}
