package test

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Unmarshal(body string) interface{} {
	var f interface{}
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		panic(err)
	}
	return f
}

func Equal(t *testing.T, got, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got = %+#v, want = %+#v", got, expected)
	}
}

func VerifyError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error returned: %s", err)
	}
}
