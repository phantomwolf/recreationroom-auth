package user

import (
	"testing"
)

func TestSetName(t *testing.T) {
	// User names to be set and expected results
	data := map[string]bool{
		"baka888": true,
		"1234":    true,
		"foo.bar": false,
		"":        false,
	}
	var user User
	for k, v := range data {
		err := user.SetName(k)
		if v && err != nil {
			t.Log("Failed to set", k, "as username:", err)
			t.Fail()
		} else if !v && err == nil {
			t.Log("Username", k, "shouldn't be set:", err)
			t.Fail()
		}
	}
}
