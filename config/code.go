package config

import (
	"../model"
)

// Constatns for role of admin
const (
	MinRoleCode   = 1000
	AdminCode     = 100
	DeveloperCode = 101
	MarketerCode  = 102
	ReportsCode   = 103
)

// Constatns for driver document
const (
	DocTypeVehicle = 0
	DocTypeDriver  = 1

	DocumentNone     = 0
	DocumentPending  = 1
	DocumentAccepted = 2

	ScheduleFromTime = 631180800
	ScheduleToTime   = 631224000
)

// Enum for Status for every models
const (
	Disabled model.Status = iota
	Enabled
	Offline = 100
	Online
	Ongoing
)

// business query constant
const (
	SortRecommend    = 1
	SortPopular      = 2
	SortDeliveryTime = 3

	QueryAll       = 0
	QueryPopular   = 1
	QueryRecommend = 2
	QueryUnder     = 3
)

// order status constant
const (
	None           = "None"
	OrderRequest   = "OrderRequest"
	OrderAccepted  = "OrderAccepted"
	OrderDeclined  = "OrderDeclined"
	OrderCancelled = "OrderCancelled"
	OrderPrepared  = "OrderPrepared"
	OrderCompleted = "OrderCompleted"
	TripRequest    = "TripRequest"
	TripAccepted   = "TripAccepted"
	TripDeclined   = "TripDeclined"
	TripCancelled  = "TripCancelled"
	TripConfirmed  = "TripConfirmed"
	TripStarted    = "TripStarted"
	TripArrived    = "TripArrived"
	TripDropped    = "TripDropped"
	TripCompleted  = "TripCompleted"
)

// const for search
const (
	DefaultSearchRadius = 8000
)

// const for reason cancel
const (
	ReasonUserCancel   = 1
	ReasonOrderDecline = 2
	ReasonTripCancel   = 3
)

// const for rate orders query
const (
	NoRated = 1
	Rated   = 2
)
