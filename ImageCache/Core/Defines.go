package Core

type BitsType int

func BitsSet(options, flag BitsType) BitsType {
	return options | flag
}

func BitsClear(options, flag BitsType) BitsType {
	return options &^ flag
}

func BitsToggle(options, flag BitsType) BitsType {
	return options ^ flag
}

func BitsHas(options, flag BitsType) bool {
	return options&flag != 0
}

type WebImageContext map[string]interface{}
