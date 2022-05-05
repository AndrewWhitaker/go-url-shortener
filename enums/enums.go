package enums

type CreationStatus int

const (
	CreationResultUnknown CreationStatus = iota
	CreationResultCreated
	CreationResultAlreadyExisted
)
