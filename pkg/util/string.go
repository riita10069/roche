package util

import "strings"

func CamelToSnake(camel string) (sneak string) {
	if camel == "" {
		return camel
	}

	delimiter := "_"
	sLen := len(camel)
	var snake string
	for i, current := range camel {
		if i > 0 && i+1 < sLen {
			if current >= 'A' && current <= 'Z' {
				next := camel[i+1]
				prev := camel[i-1]
				if (next >= 'a' && next <= 'z') || (prev >= 'a' && prev <= 'z') {
					snake += delimiter
				}
			}
		}
		snake += string(current)
	}

	snake = strings.ToLower(snake)
	return snake
}

func SnakeToUpperCamel(input string) (camel string) {
	if input == "" {
		return input
	}
	sneak := strings.Split(input, "")
	delimiter := "_"

	camel = ""
	for i := 0; i < len(sneak); i++ {
		if i == 0 {
			camel += strings.ToUpper(sneak[i])
			continue
		}

		if sneak[i] == delimiter {
			i++
			camel += strings.ToUpper(sneak[i])
			continue
		}

		camel += sneak[i]
	}
	return camel
}



func SnakeToLowerCamel(input string) (camel string) {
	if input == "" {
		return input
	}
	sneak := strings.Split(input, "")
	delimiter := "_"

	camel = ""
	for i := 0; i < len(sneak); i++ {
		if sneak[i] == delimiter {
			i++
			camel += strings.ToUpper(sneak[i])
			continue
		}
		camel += sneak[i]
	}
	return camel
}
