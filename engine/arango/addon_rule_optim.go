package arango

import (
	"github.com/gofrs/uuid"
	arngin "github.com/sudo-suhas/play-arngin"
)

type addonRuleOptim struct {
	Key uuid.UUID `json:"_key"`
	arngin.AddonRule
	Ignore ignoreClause
}

type ignoreClause struct {
	Sources      bool `json:"sources"`
	Destinations bool `json:"destinations"`
	BoardingPts  bool `json:"boardingPts"`
	DroppingPts  bool `json:"droppingPts"`
	BoardingTime bool `json:"boardingTime"`
	DroppingTime bool `json:"droppingTime"`
	SeatCount    bool `json:"seatCount"`
	BusOperators bool `json:"busOperators"`
	Duration     bool `json:"duration"`
	Appversion   bool `json:"appversion"`
	Channels     bool `json:"channels"`
}
