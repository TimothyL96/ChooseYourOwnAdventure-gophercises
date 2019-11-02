package main

import (
	"testing"
)

func TestLoadJSONFromFile(t *testing.T) {
	got, err := loadJSONFromFile()

	if err != nil {
		t.Error("error from method loadJSONFromFile:", err)
	} else {
		want := "The Little Blue Gopher1"

		if want != got["intro"].Title {
			t.Errorf("error: got %s, want %s", got["intro"].Title, want)
		}
	}

}
