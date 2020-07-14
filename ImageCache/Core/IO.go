package Core

type DataReadingOption int

const (
	ReadingMappedIfSafe DataReadingOption = iota
	ReadingUncached
	ReadingMappedAlways
)

type DataWritingOption int

const (
	WritingAtomic DataWritingOption = iota
)
