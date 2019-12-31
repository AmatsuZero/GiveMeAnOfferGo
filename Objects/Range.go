package Objects

type Range struct {
	Location uint
	Length   uint
}

func (ran *Range) IsLocationInRange(loc uint) bool {
	return !(loc < ran.Location) && (loc-ran.Location) < ran.Length
}

func (ran *Range) IsEqualTo(obj interface{}) bool {
	range2, ok := obj.(Range)
	if !ok {
		return false
	}
	return ran.Location == range2.Location
}
