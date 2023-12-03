package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	errTemplate             = fmt.Errorf("result WRONG: ")
	errNilCheck             = fmt.Errorf("failed nil check: ")
	errCompareErrorsStrings = fmt.Errorf("failed strings compare: ")
)

func CompareSameType[T comparable](received T, expected T) error {
	if received != expected {
		return fmt.Errorf("%w response: got %+v, expected %+v", errTemplate,
			received, expected)
	}

	return nil
}

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

// EqualError use check errors.Is first and then use direct string errors compare.
func EqualError(received error, expected error) error {
	if !errors.Is(received, expected) {
		if expected == nil && received != nil || received == nil && expected != nil {
			return fmt.Errorf("%w: err got: %+v err expexted: %+v", errNilCheck, received, expected)
		}

		if !(received.Error() == expected.Error()) {
			return fmt.Errorf("%w: err got %+v err expected: %+v", errCompareErrorsStrings, received, expected)
		}
	}

	return nil
}
