package main

import (
	"encoding/json"
	"testing"
)

func TestBasic(t *testing.T) {
	app := &App{}

	if err := app.Load("testdata/test.graphql"); err != nil {
		t.Fatal(err)
	}

	{
		unused := app.DetectUnused()
		if !(len(unused) == 2 && unused[0].Name == "Enum1" && unused[1].Name == "Scalar1") {
			b, _ := json.Marshal(unused)
			t.Error(string(b))
		}
	}
}
