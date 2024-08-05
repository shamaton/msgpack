package testutil

import (
	"errors"
	"maps"
	"strings"
	"testing"
)

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal(err)
	}
}

func IsError(t *testing.T, actual, expected error) {
	t.Helper()
	if !errors.Is(actual, expected) {
		t.Fatalf("not equal. actual: %v, expected: %v", actual, expected)
	}
}

func ErrorContains(t *testing.T, err error, errStr string) {
	t.Helper()
	if err == nil {
		t.Fatal("error should occur")
	}
	if !strings.Contains(err.Error(), errStr) {
		t.Fatalf("error does not contain '%s'. err: %v", errStr, err)
	}
}

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Fatalf("not equal. actual: %v, expected: %v", actual, expected)
	}
}

func EqualSlice[T comparable](t *testing.T, actual, expected []T) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("diffrent length. actual: %v, expected: %v", actual, expected)
	}
	for i := range actual {
		Equal[T](t, actual[i], expected[i])
	}
}

func EqualMap[K comparable, V comparable](t *testing.T, actual, expected map[K]V) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("diffrent length. actual: %v, expected: %v", actual, expected)
	}
	if !maps.Equal(actual, expected) {
		t.Fatalf("not equal. actual: %v, expected: %v", actual, expected)
	}
}
