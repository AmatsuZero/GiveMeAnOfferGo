package math

import (
	"GiveMeAnOffer/math/can_measure_water"
	"testing"
)

func TestCanMeasureWater(t *testing.T) {
	if !can_measure_water.CanMeasureWater(3, 5, 4) {
		t.Fail()
	}
}
