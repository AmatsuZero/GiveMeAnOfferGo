package compare_version_numbers

import (
	"strconv"
	"strings"
)

func CompareVersion(version1, version2 string) int {
	v1, v2 := strings.Split(version1, "."), strings.Split(version2, ".")
	n := len(v1)
	if n > len(v2) {
		n = len(v2)
	}

	i := 0
	for ; i < n; i++ {
		n1, _ := strconv.Atoi(trimZero(v1[i]))
		n2, _ := strconv.Atoi(trimZero(v2[i]))

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	for ; i < len(v1); i++ {
		if len(trimZero(v1[i])) > 0 {
			return 1
		}
	}

	for ; i < len(v2); i++ {
		if len(trimZero(v2[i])) > 0 {
			return -1
		}
	}

	return 0
}

func trimZero(s string) string {
	idx := 0
	for ; idx < len(s); idx++ {
		if s[idx] != '0' {
			break
		}
	}
	return s[idx:]
}
