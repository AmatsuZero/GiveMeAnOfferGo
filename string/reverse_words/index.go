package reverse_words

import "strings"

func ReverseWords(s string) string {
	var queue []string
	builder := strings.Builder{}
	for _, r := range s {
		if r != ' ' {
			builder.WriteRune(r)
		} else {
			x := builder.String()
			if x != " " {
				queue = append([]string{x}, queue...)
			}
			builder.Reset()
		}
	}
	x := builder.String()
	if x != " " {
		queue = append([]string{x}, queue...)
	}
	return strings.Join(queue, " ")
}
