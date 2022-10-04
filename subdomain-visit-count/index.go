package subdomain_visit_count

import (
	"fmt"
	"strconv"
	"strings"
)

// SubdomainVisits https://leetcode.cn/problems/subdomain-visit-count/
func SubdomainVisits(cpdomains []string) (ans []string) {
	dict := map[string]int{}

	split := func(s string) {
		i := strings.IndexByte(s, ' ')
		c, _ := strconv.Atoi(s[:i])
		s = s[i+1:]
		dict[s] += c
		for {
			i := strings.IndexByte(s, '.')
			if i < 0 {
				break
			}
			s = s[i+1:]
			dict[s] += c
		}
	}

	for _, cpdomain := range cpdomains {
		split(cpdomain)
	}

	for s, i := range dict {
		ans = append(ans, fmt.Sprintf("%v %v", i, s))
	}
	return
}
