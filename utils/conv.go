package utils

import "strings"

func ConvField(fname string) string {
	result := ""
	for i := 0; i < len(fname); i++ {
		if fname[i] == '_' {
			result += strings.ToUpper(string(fname[i+1]))
			i++
			continue
		} else {
			result += string(fname[i])
		}
	}
	return result
}

func ConvFieldFirstUpper(fname string) string {
	result := ""
	for i := 0; i < len(fname); i++ {
		if i == 0 {
			result += strings.ToUpper(string(fname[i]))
			continue
		}

		if fname[i] == '_' {
			result += strings.ToUpper(string(fname[i+1]))
			i++
			continue
		} else {
			result += string(fname[i])
		}
	}
	return result
}
