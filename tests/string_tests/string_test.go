package string_tests

import (
	"GiveMeAnOffer/leetcode"
	check_if_a_word_occurs_as_a_prefix_of_any_word_in_a_sentence "GiveMeAnOffer/string/check-if-a-word-occurs-as-a-prefix-of-any-word-in-a-sentence"
	decode_string "GiveMeAnOffer/string/decode-string"
	reformat_phone_number "GiveMeAnOffer/string/reformat-phone-number"
	"GiveMeAnOffer/string/solve_the_equation"
	string_rotation_lcci "GiveMeAnOffer/string/string-rotation-lcci"
	"GiveMeAnOffer/tree/is_palindrome"
	"testing"
)

func TestSolveEquation(t *testing.T) {
	ans := solve_the_equation.SolveEquation("x+5-3+x=6+x-2")
	if ans != "x=2" {
		t.Fail()
	}
	ans = solve_the_equation.SolveEquation("x=x")
	if ans != "Infinite solutions" {
		t.Fail()
	}
	ans = solve_the_equation.SolveEquation("2x=x")
	if ans != "x=0" {
		t.Fail()
	}
}

func TestReplaceWords(t *testing.T) {
	ans := leetcode.ReplaceWords([]string{"cat", "bat", "rat"}, "the cattle was rattled by the battery")
	if ans != "the cattle was rattled by the battery" {
		t.Fail()
	}
}

func TestMinLengthEncoding(t *testing.T) {
	ans := leetcode.MiniMumLengthEncoding([]string{"time", "me", "bell"})
	if ans != 10 {
		t.Fail()
	}
}

func TestIsPrefixOfWord(t *testing.T) {
	ans := check_if_a_word_occurs_as_a_prefix_of_any_word_in_a_sentence.IsPrefixOfWord("i love eating burger", "burg")
	if ans != 4 {
		t.Fail()
	}
	ans = check_if_a_word_occurs_as_a_prefix_of_any_word_in_a_sentence.IsPrefixOfWord("this problem is an easy problem", "pro")
	if ans != 2 {
		t.Fail()
	}
	ans = check_if_a_word_occurs_as_a_prefix_of_any_word_in_a_sentence.IsPrefixOfWord("i am tired", "you")
	if ans != -1 {
		t.Fail()
	}
}

func TestIsPalindrome(t *testing.T) {
	if !is_palindrome.IsPalindrome("A man, a plan, a canal: Panama") {
		t.Fail()
	}

	if is_palindrome.IsPalindrome("race a car") {
		t.Fail()
	}

	if !is_palindrome.IsPalindrome(".,") {
		t.Fail()
	}

	if is_palindrome.IsPalindrome("0P") {
		t.Fail()
	}
}

func TestDecodeString(t *testing.T) {
	if decode_string.DecodeString("3[a]2[bc]") != "aaabcbc" {
		t.Fail()
	}
}

func TestIsFlippedString(t *testing.T) {
	if !string_rotation_lcci.IsFlippedString("waterbottle", "erbottlewat") {
		t.Fail()
	}

	if !string_rotation_lcci.IsFlippedString("", "") {
		t.Fail()
	}
}

func TestReformatPhoneNumber(t *testing.T) {
	if reformat_phone_number.ReformatNumber("1-23-45 6") != "123-456" {
		t.Fail()
	}

	if reformat_phone_number.ReformatNumber("123 4-567") != "123-45-67" {
		t.Fail()
	}

	if reformat_phone_number.ReformatNumber("123 4-5678") != "123-456-78" {
		t.Fail()
	}
}
