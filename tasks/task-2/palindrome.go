package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	str1 := "noon"
	fmt.Println(is_palindrome(str1))

	str2 := "hello world"
	fmt.Println(is_palindrome(str2))
}

func is_palindrome(str string) bool {
	lowerStr := strings.ToLower(str)

	reg := regexp.MustCompile(`[^\w]`)
	cleanedStr := reg.ReplaceAllString(lowerStr, "")

	// palindrome check
	left, right := 0, len(cleanedStr)-1
	for left < right {
		if cleanedStr[left] != cleanedStr[right] {
			return false
		}
		left++
		right--
	}
	return true
}
