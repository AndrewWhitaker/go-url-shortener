package enums

type CreationStatus int

const (
	CreationResultUnknown CreationStatus = iota
	CreationResultCreated
	CreationResultAlreadyExists
	CreationResultDuplicateSlug
	CreationResultUnknownError
)
