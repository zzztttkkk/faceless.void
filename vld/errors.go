package vld

type ErrorKind int

const (
	ErrorKindMissingRequired = ErrorKind(iota)

	ErrorKindCustomFunc

	ErrorKindIntLtMin
	ErrorKindIntGtMax
	ErrorKindIntNotInRange

	ErrorKindStringTooLong
	ErrorKindStringTooShort
	ErrorKindStringNotMatched
	ErrorKindStringNotInRanges

	ErrorKindContainerSizeTooLarge
	ErrorKindContainerSizeTooSmall
)
