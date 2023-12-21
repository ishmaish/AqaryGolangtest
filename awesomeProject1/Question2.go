package main

import (
	"fmt"
	"sort"
)

func rearr(s string) string {

	// store the number of rune
	charCountMap := make(map[rune]int)

	for _, char := range s {
		charCountMap[char]++
	}

	// to storing chars purpose
	sortedChars := make([]rune, 0, len(charCountMap))

	for char := range charCountMap {
		sortedChars = append(sortedChars, char)
	}

	sort.Slice(sortedChars, func(i, j int) bool {
		return charCountMap[sortedChars[i]] > charCountMap[sortedChars[j]]
	})

	result := make([]rune, len(s))

	index := 0
	for _, char := range sortedChars {
		for charCountMap[char] > 0 && index < len(s) {
			result[index] = char
			index = index + 2
			charCountMap[char]--
		}
	}

	index = 1
	for _, char := range sortedChars {
		for charCountMap[char] > 0 && index < len(s) {
			result[index] = char
			index = index + 2
			charCountMap[char]--
		}
	}

	for _, char := range sortedChars {
		if charCountMap[char] > 0 {
			return " "
		}
	}

	for i := 0; i < len(result); i++ {
		for j := 1; j < len(result); j++ {
			if result[i] != result[j] {
				return string(result)
			}
		}
	}

	return ""
}

func main() {

	s := "ccc"
	res1 := rearr(s)
	fmt.Println(res1)

	s1 := "ccb"
	res2 := rearr(s1)
	fmt.Println(res2)
}
