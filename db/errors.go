package db

import "strings"

func CheckIfDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "#23505")
}

func CheckIfShardExists(err error) bool {
	return strings.Contains(err.Error(), "#42P01")
}
