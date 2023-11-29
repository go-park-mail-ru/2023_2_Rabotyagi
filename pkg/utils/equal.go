package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var errTemplate = fmt.Errorf("result WRONG: ")

func EqualTest(received any, expected any) error {
	if !reflect.DeepEqual(received, expected) {
		expectedRaw, err := json.Marshal(expected)
		if err != nil {
			return err //nolint:wrapcheck
		}

		return fmt.Errorf("%w response: got %s, expected %s", errTemplate,
			string(expectedRaw), string(expectedRaw))
	}

	return nil
}
