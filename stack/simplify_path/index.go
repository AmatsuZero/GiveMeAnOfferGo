package simplifypath

import "strings"

// https://leetcode-cn.com/problems/simplify-path/
func SimplifyPath(path string) string {
	var ret []string
	components := strings.Split(path, "/")
	for _, com := range components {
		switch com {
		case "..":
			if len(ret) > 0 {
				ret = ret[:len(ret)-1]
			}
		default:
			if len(com) > 0 && com != "." {
				ret = append(ret, com)
			}
		}
	}
	if len(ret) == 0 {
		return "/"
	}
	return "/" + strings.Join(ret, "/")
}
