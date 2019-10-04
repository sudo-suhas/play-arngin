package opa

import (
	"bytes"
	"strconv"

	arngin "github.com/sudo-suhas/play-arngin"
)

func toOPARule(r arngin.AddonRule) []byte {
	var buf bytes.Buffer

	buf.WriteRune('{')
	writeIntSetRule(&buf, "sources", r.Sources)
	writeIntSetRule(&buf, "destinations", r.Destinations)
	writeIntSetRule(&buf, "boardingPts", r.BoardingPts)
	writeIntSetRule(&buf, "droppingPts", r.DroppingPts)
	writeFloatRangeRule(&buf, "boardingTime", r.BoardingTime)
	writeFloatRangeRule(&buf, "droppingTime", r.DroppingTime)
	writeNumCompareRule(&buf, "seatCount", r.SeatCount)
	writeStringSetRule(&buf, "busOperators", r.BusOperators)
	writeNumCompareRule(&buf, "duration", r.Duration)
	writeNumCompareRule(&buf, "appversion", r.Appversion)
	writeStringSetRule(&buf, "channels", r.Channels)
	writeStringSetRule(&buf, "addonIDs", r.Addons)
	buf.WriteRune('}')

	return buf.Bytes()
}
func writeIntSetRule(buf *bytes.Buffer, field string, a []int) {
	if len(a) == 0 {
		return
	}

	buf.WriteString(`"` + field + `": {`)
	buf.Write(intsCSV(a))
	buf.WriteString("}, ")
}

func intsCSV(a []int) []byte {
	// Appr. 3 chars per num plus the comma.
	estimate := len(a) * 4
	b := make([]byte, 0, estimate)
	for _, n := range a {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1] // strip last comma
	return b
}

func writeFloatRangeRule(buf *bytes.Buffer, field string, fr *arngin.FloatRange) {
	if fr == nil || fr.GTE == 0 || fr.LTE == 0 {
		return
	}

	buf.WriteString(`"` + field + `": [`)
	buf.Write(strconv.AppendFloat(nil, fr.GTE, 'f', 3, 64))
	buf.WriteRune(',')
	buf.Write(strconv.AppendFloat(nil, fr.LTE, 'f', 3, 64))
	buf.WriteString("], ")
}

func writeStringSetRule(buf *bytes.Buffer, field string, a []string) {
	if len(a) == 0 {
		return
	}

	buf.WriteString(`"` + field + `": {`)
	for i, s := range a {
		if i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(`"` + s + `"`)
	}
	buf.WriteString("}, ")
}

func writeNumCompareRule(buf *bytes.Buffer, field string, r *arngin.NumRule) {
	if r == nil {
		return
	}

	buf.WriteString(`"` + field + `": { "op": "` + string(r.Op) + `", "value": `)
	buf.Write(strconv.AppendInt(nil, int64(r.Value), 10))
	buf.WriteString("}, ")
}
