package netcat

import "regexp"

func IsValidName(name string) bool {
	//regex for alphanumeric characters only
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	//return true if length is > 0 and <= 20 and all characters are alphanumeric
	return len(name) > 0 && len(name) <= 20 && (regex.MatchString(name))
}
