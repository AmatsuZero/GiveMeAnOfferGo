package array

import "testing"

func TestSearchInsert(t *testing.T) {
	input := []int{1, 2, 3, 5, 6}

	if SearchInsert(input, 5) == 2 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if SearchInsert(input, 2) == 1 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if SearchInsert(input, 7) == 4 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if SearchInsert(input, 0) == 0 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}

func TestSnippets(t *testing.T) {
	input := []int{1, 2, 3, 5, 6}
	_, input = input[len(input)-1], input[:len(input)-1]
	_, input = input[0], input[1:]
	t.Log(input)
}
