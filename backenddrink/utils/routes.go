package utils

import (
	"errors"
	"regexp"
)

func ValidateUserPath(path string) error {
	matched, err := regexp.Match(
		`(https:\/\/cdn\.discordapp\.com\/attachments\/){1}(\w|\W)*(\.png|\.jpeg|\.jpg|\.gif)$`,
		[]byte(path),
	)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid path")
	}
	return nil
}
