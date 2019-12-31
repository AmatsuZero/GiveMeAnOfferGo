package queue

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/test/Utils"
	"fmt"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	buffer := Collections.NewBufferRing(5)
	getInt := Utils.GetInt
	buffer.Write(getInt(123))
	buffer.Write(getInt(456))
	buffer.Write(getInt(789))
	buffer.Write(getInt(666))

	fmt.Println(buffer.Read())
	fmt.Println(buffer.Read())
	fmt.Println(buffer.Read())

	buffer.Write(getInt(333))
	buffer.Write(getInt(555))

	fmt.Println(buffer.Read())
	fmt.Println(buffer.Read())
	fmt.Println(buffer.Read())
	fmt.Println(buffer.Read())
}
