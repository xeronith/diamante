package search

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// TODO: Performance drops by a factor of 10 with normalization.
func MatchAny(input string, criteria string) bool {
	if criteria == "" {
		return true
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	//replacer := strings.NewReplacer("ي","ة","ك","ی","ه","ک")
	normalizedCriteria, _, err := transform.String(t, criteria)
	//normalizedCriteriaReplaced := replacer.Replace(normalizedCriteria)
	normalizedCriteria = strings.ReplaceAll(normalizedCriteria, "ي", "ی")
	normalizedCriteria = strings.ReplaceAll(normalizedCriteria, "ك", "ک")
	normalizedCriteria = strings.ReplaceAll(normalizedCriteria, "ة", "ه")

	if err != nil {
		//TODO: Handle
		return false
	}

	normalizedInput, _, err := transform.String(t, input)
	//normalizedInputReplaced := replacer.Replace(normalizedInput)
	normalizedInput = strings.ReplaceAll(normalizedInput, "ي", "ی")
	normalizedInput = strings.ReplaceAll(normalizedInput, "ك", "ک")
	normalizedInput = strings.ReplaceAll(normalizedInput, "ة", "ه")

	if err != nil {
		//TODO: Handle
		return false
	}

	input = strings.ToUpper(normalizedInput)
	criteria = strings.ToUpper(normalizedCriteria)

	terms := strings.Split(criteria, " ")
	for _, term := range terms {
		trimmedTerm := strings.TrimSpace(term)
		if trimmedTerm != "" && strings.Contains(input, trimmedTerm) {
			return true
		}
	}

	return false
}

func CountRepeatedChars(word string, searchWord string) int {
	var sum = 0
	for _, term := range strings.Split(searchWord, "") {
		if strings.Contains(word, term) {
			sum += 1
		}
	}

	for _, item := range strings.Split(word, " ") {
		if strings.HasPrefix(searchWord, item) {
			sum += 5
		}
	}

	if word == searchWord {
		sum += 100
	}

	return sum
}
