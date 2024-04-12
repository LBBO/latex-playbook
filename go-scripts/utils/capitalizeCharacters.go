package utils

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func CapitalizeCharacters(characters []string) (capitalizedCharacters []string, err error) {
	for _, character := range characters {
		r, size := utf8.DecodeRuneInString(character)
		if r == utf8.RuneError {
			return nil, fmt.Errorf("could not decode rune in string: %+v", r)
		}
		capitalizedCharacters = append(capitalizedCharacters, string(unicode.ToUpper(r))+character[size:])
	}

	return
}
