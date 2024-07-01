package constant

const (
	ServiceName        = "tonfura-exercise"
	ResponseCodePrefix = 1

	StatusAvailable = "available"
	StatusSoldOut   = "sold_out"

	BookingStatusToBeConfirm = "to_be_confirmed"
	BookingStatusInProgress  = "in_progress"
	BookingStatusCheckedIn   = "checked_in"

	BookingCheckInSuccess = "success"
	BookingCheckInFail    = "fail"

	ClassTypeEconomy  = "economy"
	ClassTypeBusiness = "business"
	ClassTypeFirst    = "first"
)

const (
	ClassCodeEconomy int = iota
	ClassCodeBusiness
	ClassCodeFirst
)

var ClassTypes = map[string]int{
	ClassTypeEconomy:  ClassCodeEconomy,
	ClassTypeBusiness: ClassCodeBusiness,
	ClassTypeFirst:    ClassCodeFirst,
}
