package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	text := "Hello world from Go! This is GO. Hello World, again."
	frequency := wordFrequency(text)
	fmt.Println(frequency)
}

func wordFrequency(inputText string) map[string]int {
	lowerText := strings.ToLower(inputText)

	reg := regexp.MustCompile(`[^\w\s]`)
	cleanText := reg.ReplaceAllString(lowerText, "")

	words := strings.Fields(cleanText)
	fmt.Println(words)

	freq_map := make(map[string]int)
	for _, word := range words {
		freq_map[word]++
	}

	return freq_map
}
