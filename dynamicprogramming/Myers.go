package dynamicprogramming

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/**
https://github.com/cj1128/myers-diff
*/
type operation uint

const (
	INSERT operation = iota
	DELETE
	MOVE
)

func (op operation) String() string {
	switch op {
	case INSERT:
		return "INS"
	case DELETE:
		return "DEL"
	case MOVE:
		return "MOV"
	default:
		return "UNKNOWN"
	}
}

var colors = map[operation]string{
	INSERT: "\033[32m",
	DELETE: "\033[31m",
	MOVE:   "\033[39m",
}

func generateDiff(src, dst []string) (ret string) {
	script := shortestEditScript(src, dst)
	srcIndex, dstIndex := 0, 0
	for _, op := range script {
		switch op {
		case INSERT:
			ret += fmt.Sprintf(colors[op] + "+" + dst[dstIndex])
			dstIndex += 1
		case MOVE:
			ret += fmt.Sprintf(colors[op] + " " + src[srcIndex])
			srcIndex += 1
			dstIndex += 1
		case DELETE:
			ret += fmt.Sprintf(colors[op] + "-" + src[srcIndex])
			srcIndex += 1
		}
		ret += "\n"
	}
	return
}

// 生成最短的编辑脚本
func shortestEditScript(src, dst []string) []operation {
	n, m := len(src), len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		// 最多只有d+1个k
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		// 需要注意处理对角线
		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			// 向下
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else { // 向右
				x = lastV[k-1] + 1
			}

			y = x - k

			// 处理对角线
			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	// for debug
	printTrace(trace)

	// 反向回溯
	var script []operation

	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, MOVE)
			x--
			y--
		}

		if x == prevX {
			script = append(script, INSERT)
		} else {
			script = append(script, DELETE)
		}

		x, y = prevX, prevY
	}

	if len(trace) > 0 && len(trace[0]) > 0 && trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, MOVE)
		}
	}

	for left, right := 0, len(script)-1; left < right; left, right = left+1, right-1 {
		script[left], script[right] = script[right], script[left]
	}
	return script
}

func FileDiff(input string, output string) string {
	src, err := getFileLines(input)
	if err != nil {
		log.Fatal(err)
	}
	dst, err := getFileLines(output)
	if err != nil {
		log.Fatal(err)
	}
	return generateDiff(src, dst)
}

func getFileLines(p string) ([]string, error) {
	f, err := os.Open(p)

	if err != nil {
		return nil, err
	}

	defer func() {
		err2 := f.Close()
		if err2 != nil {
			fmt.Println(err2)
		}
	}()

	scanner := bufio.NewScanner(f)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func printTrace(trace []map[int]int) {
	for d := 0; d < len(trace); d++ {
		fmt.Printf("d = %d:\n", d)
		v := trace[d]
		for k := -d; k <= d; k += 2 {
			x := v[k]
			y := x - k
			fmt.Printf("  k = %2d: (%d, %d)\n", k, x, y)
		}
	}
}
