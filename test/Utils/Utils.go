package Utils

import "github.com/AmatsuZero/GiveMeAnOfferGo/Objects"

func GetInt(num int) *Objects.NumberObject {
	return Objects.NewNumberWithInt(num)
}

func GetString(str string) *Objects.StringObject {
	return &Objects.StringObject{GoString: &str}
}
