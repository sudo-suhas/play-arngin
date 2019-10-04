package arngin

import data.arngin.rules

# Match any and all criteria specified in each rule against the input.
# Extract the list of addon IDs from rules which do match into the set
# matched_addons which can be accessed after evaluating the query as
# data.arngin.matched_addons.
matched_addons[rule.addonIDs[_]] {
	rule := rules[_]
	match_source(rule, input.source)
	match_destination(rule, input.destination)
	match_boarding_time(rule, input.boardingTime)
	match_dropping_time(rule, input.droppingTime)
	match_seats(rule, input.seats)
	match_bus_operator(rule, input.busOperator)
	match_duration(rule, input.duration)
	match_appversion(rule, input.appversion)
	match_channel(rule, input.channel)
}

# ------------ Start source match function  ------------

match_source(r, v) {
	not r.sources
}

match_source(r, v) {
	r.sources[v]
}

# ------------- End source match function  -------------
# ------------------------------------------------------
# ---------- Start destination match function  ---------

match_destination(r, v) {
	not r.destinations
}

match_destination(r, v) {
	r.destinations[v]
}

# ----------- End destination match function  ----------
# ------------------------------------------------------
# --------- Start boarding time match function  --------

match_boarding_time(r, v) {
	not r.boardingTime
}

match_boarding_time(r, v) {
	r.boardingTime[0] <= v
	r.boardingTime[1] >= v
}

# ---------- End boarding time match function  ---------
# ------------------------------------------------------
# --------- Start dropping time match function  --------

match_dropping_time(r, v) {
	not r.droppingTime
}

match_dropping_time(r, v) {
	r.droppingTime[0] <= v
	r.droppingTime[1] >= v
}

# ---------- End dropping time match function  ---------
# ------------------------------------------------------
# ---------- Start seat count match function  ----------

match_seats(r, v) {
	not r.seatCount
}

match_seats(r, v) {
	r.seatCount.op == "eq"
	r.seatCount.value == v
}

match_seats(r, v) {
	r.seatCount.op == "lte"
	r.seatCount.value <= v
}

match_seats(r, v) {
	r.seatCount.op == "gte"
	r.seatCount.value >= v
}

# ----------- End seat count match function  -----------
# ------------------------------------------------------
# --------- Start bus operator match function  ---------

match_bus_operator(r, v) {
	not r.busOperators
}

match_bus_operator(r, v) {
	r.busOperators[v]
}

# ---------- End bus operator match function  ----------
# ------------------------------------------------------
# -------- Start travel duration match function  -------

match_duration(r, v) {
	not r.duration
}

match_duration(r, v) {
	r.duration.op == "eq"
	r.duration.value == v
}

match_duration(r, v) {
	r.duration.op == "lte"
	r.duration.value <= v
}

match_duration(r, v) {
	r.duration.op == "gte"
	r.duration.value >= v
}

# --------- End travel duration match function  --------
# ------------------------------------------------------
# ---------- Start app version match function  ---------

match_appversion(r, v) {
	not r.appversion
}

match_appversion(r, v) {
	v == 0 # we didn't get a version to check
}

match_appversion(r, v) {
	r.appversion.op == "eq"
	r.appversion.value == v
}

match_appversion(r, v) {
	r.appversion.op == "lte"
	r.appversion.value <= v
}

match_appversion(r, v) {
	r.appversion.op == "gte"
	r.appversion.value >= v
}

# ----------- End app version match function  ----------
# ------------------------------------------------------
# --------- Start sales channel match function  --------

match_channel(r, v) {
	not r.channels
}

match_channel(r, v) {
	r.channels[v]
}

# ----------- End sales channel match function  ---------
