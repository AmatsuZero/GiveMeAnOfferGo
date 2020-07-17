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

type WebImageContext map[string]interface{}
