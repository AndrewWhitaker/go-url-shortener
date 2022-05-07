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
