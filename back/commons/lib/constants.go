package lib

import (
	"app/commons/constants"
	"errors"
)

func IsCustomError(err error) bool {
	for _, customErr := range constants.CUSTOM_ERRORS {
		if errors.Is(customErr.Err, err) {
			return true
		}
	}

	return false
}
