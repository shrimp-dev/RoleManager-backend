package utils

import (
	"errors"
	"fmt"
	"reflect"
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

func ValidateBody(body interface{}, ignore ...string) bool {
	v := reflect.ValueOf(body)
	f := []string{}
	if len(f) == 0 {
		for i := 0; i < v.NumField(); i++ {
			f = append(f, v.Type().Field(i).Name)
		}

	}

	for i := 0; i < len(f); i++ {
		if found := findAndRemoveField(&ignore, f[i]); found {
			continue
		}

		if isZero := v.FieldByName(f[i]).IsZero(); isZero {
			return false
		}

	}
	return true
}

func findAndRemoveField(f *[]string, name string) bool {
	for i := 0; i < len(*f); i++ {
		if (*f)[i] == name {
			fmt.Println(name)
			*f = append((*f)[:i], (*f)[i+1:]...)
			return true
		}
	}
	return false
}
