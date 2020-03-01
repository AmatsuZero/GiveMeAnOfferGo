package math

import (
	"GiveMeAnOfferGo/math"
	"testing"
)

func TestAddBinary(t *testing.T) {
	if math.AddBinary("11", "1") != "100" {
		t.Fail()
	}

	if math.AddBinary("1010", "1011") != "10101" {
		t.Fail()
	}

	if math.AddBinary("10100000100100110110010000010101111011011001101110111111111101000000101111001110001111100001101",
		"110101001011101110001111100110001010100001101011101010000011011011001011101111001100000011011110011") != "110111101100010011000101110110100000011101000101011001000011011000001100011110011010010011000000000" {
		t.Fail()
	}
}
