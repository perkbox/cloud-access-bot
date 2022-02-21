package utils

import "strings"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SplitFreeString(inputString string) []string {
	var cleanedValues []string
	//----Splits by new line and commas
	splitNewLine := strings.Split(inputString, "\n")
	for _, v := range splitNewLine {
		splitComma := strings.Split(v, ",")
		cleanedValues = append(cleanedValues, splitComma...)
	}
	return cleanedValues
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
