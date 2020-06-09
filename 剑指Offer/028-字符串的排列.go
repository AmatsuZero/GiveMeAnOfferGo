package 剑指Offer

import (
	"fmt"
)

/*
题目：输入一个字符串，打印出该字符串中字符的所有排列。
例如输入字符串abc，则打印出由字符a、b、c所能排列出来的所有字符串abc、acb、bac、bca、cab和cba。
*/

func Permutation(str string) {
	chars := []rune(str)
	permutation(0, &chars)
}

/*
start 指向排列字符串的第一个
*/
func permutation(start int, arranged *[]rune) {
	if len(*arranged) == start {
		fmt.Println(string(*arranged))
		return
	}
	for i := 0; i < len(*arranged); i++ {
		if i != start {
			(*arranged)[start], (*arranged)[i] = (*arranged)[i], (*arranged)[start]
		}
		permutation(start+1, arranged)
		if i != start {
			(*arranged)[start], (*arranged)[i] = (*arranged)[i], (*arranged)[start]
		}
	}
}

/*
   输入一个含有8个数字的数组，判断有么有可能把这8个数字分别放到正方体的8个顶点上，使得正方体上三组相对的面上的4个顶点的和相等。
   思路：相当于求出8个数字的全排列，判断有没有一个排列符合题目给定的条件，即三组对面上顶点的和相等。
*/
func CubVertex(arranged [8]int) (result bool) {
	return cubVertex(0, &arranged)
}

func cubVertex(begin int, arranged *[8]int) (result bool) {
	if arranged == nil {
		return
	}
	/*
		这相当于先得到a1、a2、a3、a4、a5、a6、a7和a8这8个数字的所有排列，然后判断有没有某一个的排列符合题目给定的条件，
		即a1＋a2＋a3＋a4＝a5＋a6＋a7＋a8，a1＋a3＋a5＋a7＝a2＋a4＋a6＋a8，并且a1＋a2＋a5＋a6＝a3＋a4＋a7＋a8
	*/
	if begin == len(*arranged)-1 {
		result = (*arranged)[0]+(*arranged)[1]+(*arranged)[2]+(*arranged)[3] == (*arranged)[4]+(*arranged)[5]+(*arranged)[6]+(*arranged)[7] &&
			(*arranged)[0]+(*arranged)[2]+(*arranged)[4]+(*arranged)[6] == (*arranged)[1]+(*arranged)[3]+(*arranged)[5]+(*arranged)[7] &&
			(*arranged)[0]+(*arranged)[1]+(*arranged)[4]+(*arranged)[5] == (*arranged)[2]+(*arranged)[3]+(*arranged)[6]+(*arranged)[7]
	} else {
		for i := begin; i < len(*arranged); i++ {
			if i != begin {
				(*arranged)[begin], (*arranged)[i] = (*arranged)[i], (*arranged)[begin]
			}
			result = cubVertex(begin+1, arranged)
			if result {
				break
			}
			if i != begin {
				(*arranged)[begin], (*arranged)[i] = (*arranged)[i], (*arranged)[begin]
			}
		}
	}
	return
}

type chessBoard struct {
	//初始值默认都为0,为便于扩展，空间大小乘以了10
	//pos[i]=j表示第i行皇后放在j位置
	//用b[i]表示第i行是否摆放了皇后ba取值在0~7）
	pos [8]int
	b   [8]bool
	//用c[j-i+7]正对角线是否摆放了皇后（c取值在0~14）
	//用d[i+j]斜对角线是否摆法了皇后（d取值在0~14）
	c, d [15]bool
}

// 落子
func (board *chessBoard) put(i, j int, n bool) {
	board.pos[i], board.b[j], board.c[j-i+7], board.d[i+j] = j, n, n, n
}

/*检查相应位置是否可以放置皇后*/
func (board *chessBoard) checkPos(i, j int) bool {
	return !(board.b[j] || board.c[j-i+7] || board.d[i+j])
}

func (board *chessBoard) String() string {
	str := ""
	for i := 0; i < len(board.pos); i++ {
		for j := 0; j < len(board.pos); j++ {
			if board.pos[i] == j {
				str += "1"
			} else {
				str += "0"
			}
		}
		str += "\n"
	}
	str += "=================="
	return str
}

/*
在8×8的国际象棋上摆放8个皇后，使其不能相互攻击，即任意两个皇后不得处在同一行、同一列或者同一对角线上。请问总共有多少种符合条件的摆法
*/
func Queen() (count int) {
	board := chessBoard{}
	board.queenResolve(0, &count)
	return
}

func (board *chessBoard) queenResolve(i int, count *int) {
	if i > 7 { // 说明已经获得了解
		*count++
		fmt.Println(board)
		return
	}
	for j := 0; j < len(board.pos); j++ {
		if board.checkPos(i, j) {
			board.put(i, j, true)
			board.queenResolve(i+1, count)
			board.put(i, j, false)
		}
	}
}
