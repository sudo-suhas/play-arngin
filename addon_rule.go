package arngin

import "github.com/gofrs/uuid"

// AddonRule is composed of optional rule attributes and the list of
// addon IDs.
type AddonRule struct {
	ID           uuid.UUID   `json:"id"`
	Sources      []int       `json:"sources,omitempty"`
	Destinations []int       `json:"destinations,omitempty"`
	BoardingPts  []int       `json:"boardingPts,omitempty"`
	DroppingPts  []int       `json:"droppingPts,omitempty"`
	BoardingTime *FloatRange `json:"boardingTime,omitempty"`
	DroppingTime *FloatRange `json:"droppingTime,omitempty"`
	SeatCount    *NumRule    `json:"seatCount,omitempty"`
	BusOperators []string    `json:"busOperators,omitempty"`
	Duration     *NumRule    `json:"duration,omitempty"`
	Appversion   *NumRule    `json:"appversion,omitempty"`
	Channels     []string    `json:"channels,omitempty"`
	Addons       []string    `json:"addons"`
}

// NumRule defines the relation and value to compare against for the
// input.
type NumRule struct {
	Op    CompareOp `json:"op" validate:"required,oneof=eq lte gte"`
	Value int       `json:"value" validate:"required,gte=0"`
}

// FloatRange is a range of floating point numbers that the input is
// expected to match against.
type FloatRange struct {
	GTE float64 `json:"gte"`
	LTE float64 `json:"lte"`
}

// CompareOp is one of equal to, greater than equal to and less than or
// equal to.
type CompareOp string

// Compare operators
const (
	EqOp  CompareOp = "eq"
	GTEOp CompareOp = "gte"
	LTEOp CompareOp = "lte"
)
