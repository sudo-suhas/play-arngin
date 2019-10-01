package arngin

// AddonsQ is the input which needs to be matched against addon rules.
type AddonsQ struct {
	Source       int     `json:"source"`
	Destination  int     `json:"destination"`
	BoardingPt   int     `json:"boardingPt"`
	DroppingPt   int     `json:"droppingPt"`
	BoardingTime float64 `json:"boardingTime"`
	DroppingTime float64 `json:"droppingTime"`
	Seats        int     `json:"seats"`
	BusOperator  string  `json:"busOperator"`
	Duration     int     `json:"duration"`
	Appversion   int     `json:"appversion"`
	Channel      string  `json:"channel"` // One of WEB_DIRECT, MOBILE_APP, MOBILE_WEB
}
