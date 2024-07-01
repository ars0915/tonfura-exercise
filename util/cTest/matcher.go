package ctest

import (
	"reflect"
	"slices"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

// reference: https://github.com/golang/mock/issues/616#issuecomment-1040155899
type mockDiffWrapper struct {
	v interface{}
}

func DiffWrapper(v interface{}) *mockDiffWrapper {
	return &mockDiffWrapper{
		v: v,
	}
}

// Implement the gomock.Matcher interface
func (m *mockDiffWrapper) Matches(v interface{}) bool {
	ok := gomock.Eq(m.v).Matches(v) // wrap the gomock.Eq matcher
	if !ok {
		logrus.Errorf("This is the diff want: %s", pretty.Diff(m.v, v))
	}
	return ok
}

// Implement the gomock.Matcher interface
func (m *mockDiffWrapper) String() string {
	return gomock.Eq(m.v).String() // again wrap the gomock.Eq matcher
}

type exceptMatcher struct {
	v      interface{}
	fields []string
}

func ExceptMatcher(v interface{}, fields []string) *exceptMatcher {
	return &exceptMatcher{v, fields}
}

func (m *exceptMatcher) Matches(x interface{}) (result bool) {
	defer func() {
		if !result {
			logrus.Errorf("This is the diff want: %s", pretty.Diff(m.v, x))
		}
	}()

	expectValue := reflect.Indirect(reflect.ValueOf(m.v))
	actualValue := reflect.Indirect(reflect.ValueOf(x))

	for i := 0; i < expectValue.NumField(); i++ {
		typeField := expectValue.Type().Field(i)

		if slices.Contains(m.fields, typeField.Name) {
			continue
		}

		fieldExpect := expectValue.Field(i)
		fieldActual := actualValue.FieldByName(typeField.Name)

		if result = checkFieldEqual(fieldExpect, fieldActual); !result {
			return
		}
	}

	result = true
	return
}

func (m *exceptMatcher) String() string {
	return gomock.Eq(m.v).String()
}

func checkFieldEqual(expected, actual reflect.Value) bool {
	if !expected.IsValid() || !actual.IsValid() {
		return false
	}

	expectedType := expected.Type()
	actualType := actual.Type()
	if expectedType != actualType {
		return false
	}

	switch expectedType.Kind() {
	case reflect.Interface:
		if ok := checkFieldEqual(expected.Elem(), actual.Elem()); !ok {
			return false
		}
	case reflect.Map:
		expectedKeys := expected.MapKeys()
		actualKeys := actual.MapKeys()

		if len(expectedKeys) != len(actualKeys) {
			return false
		}

		for _, i := range expectedKeys {
			if ok := checkFieldEqual(expected.MapIndex(i), actual.MapIndex(i)); !ok {
				return false
			}
		}
	case reflect.Ptr:
		switch {
		case expected.IsNil() && !actual.IsNil():
			return false
		case !expected.IsNil() && actual.IsNil():
			return false
		case !expected.IsNil() && !actual.IsNil():
			if ok := checkFieldEqual(expected.Elem(), actual.Elem()); !ok {
				return false
			}
		}
	case reflect.Array, reflect.Slice:
		lenA := expected.Len()
		lenB := actual.Len()
		if lenA != lenB {
			return false
		}
		for i := 0; i < lenA; i++ {
			if ok := checkFieldEqual(expected.Index(i), actual.Index(i)); !ok {
				return false
			}
		}
	case reflect.Struct:
		switch e := expected.Interface().(type) {
		case time.Time:
			a, ok := actual.Interface().(time.Time)
			if !ok {
				return false
			}

			if e.Unix() != a.Unix() {
				return false
			}

			return true
		}
		for i := 0; i < expected.NumField(); i++ {
			if ok := checkFieldEqual(expected.Field(i), actual.Field(i)); !ok {
				return false
			}
		}
	default:
		if !expected.CanInterface() || !actual.CanInterface() {
			return false
		}
		if expected.Interface() != actual.Interface() {
			return false
		}
	}

	return true
}
