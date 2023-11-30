package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var errTemplate = fmt.Errorf("result WRONG: ")

// EqualTest accept arguments that should be work with json.Marshal
func EqualTest(received any, expected any) error {
	if !reflect.DeepEqual(received, expected) {
		expectedRaw, err := json.Marshal(expected)
		if err != nil {
			return err //nolint:wrapcheck
		}

		receivedRaw, err := json.Marshal(expected)
		if err != nil {
			return err //nolint:wrapcheck
		}

		return fmt.Errorf("%w response: got %s, expected %s", errTemplate,
			string(receivedRaw), string(expectedRaw))
	}

	return nil
}