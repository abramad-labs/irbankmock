package security

import "regexp"

func StringHasInsecureCharacters(content string) bool {
	matched, err := regexp.MatchString("[%&$#`<>'\"\\{}=]+", content)
	if err != nil {
		return true
	}
	return matched
}
