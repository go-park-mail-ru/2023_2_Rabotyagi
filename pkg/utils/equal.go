package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	errTemplate      = fmt.Errorf("result WRONG: ")
	errCompareErrors = fmt.Errorf("failed err compare: ")
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

		receivedRaw, err := json.Marshal(received)
		if err != nil {
			return err //nolint:wrapcheck
		}

		return fmt.Errorf("%w response: got %s, expected %s", errTemplate,
			string(receivedRaw), string(expectedRaw))
	}

	return nil
}

func EqualError(received error, expected error) error {
	if !errors.Is(received, expected) {
		return fmt.Errorf("%w got %+v expected wrapped: %+v", errCompareErrors, received, expected)
	}

	return nil
}
