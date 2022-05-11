package enums

type CreationStatus int

const (
	CreationResultUnknown CreationStatus = iota
	CreationResultCreated
	CreationResultAlreadyExists
	CreationResultDuplicateSlug
	CreationResultUnknownError
)

type DeleteStatus int

const (
	DeleteResultUnknown DeleteStatus = iota
	DeleteResultSuccessful
	DeleteResultNotFound
	DeleteResultUnknownError
)

type GetClicksStatus int

const (
	GetClicksResultUnknown GetClicksStatus = iota
	GetClicksResultSuccessful
	GetClicksResultNotFound
	GetClicksResultUnknownError
)

type GetClicksTimePeriod int

const (
	GetClicksTimePeriodAllTime GetClicksTimePeriod = iota
	GetClicksTimePeriodPastWeek
	GetClicksTimePeriod24Hours
)
