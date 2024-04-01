package main

type ModeEnum int

const (
	ModeOne ModeEnum = iota
	ModeTwo
	ModeThree
)

var ModeEnumValues = []string{
	"ModeOne",
	"ModeTwo",
	"ModeThree",
}

func (e ModeEnum) String() string {
	return ModeEnumValues[e]
}
