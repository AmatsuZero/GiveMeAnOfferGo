package 剑指Offer

import (
	. "github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"math/rand"
)

func RandomIntArray(len int) (output []int) {
	if len <= 0 {
		return
	}
	output = make([]int, len)
	for i := 0; i <= len-1; i++ {
		output[i] = rand.Intn(len)
	}
	return output
}

func RandomIntList(length int) (pHead *ListNode) {
	source := RandomIntArray(length)
	if len(source) == 0 {
		return
	}
	pHead = &ListNode{Val: source[0]}
	if len(source) == 1 {
		return
	}
	for i := 1; i < len(source); i++ {
		pHead.AddToTail(source[i])
	}
	return
}
