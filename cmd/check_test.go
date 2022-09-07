package cmd

import (
	"testing"
)

func TestTagsVar(t *testing.T) {
	str := "lorem=ipsum;dolor=sit-amet"
	tags := tagsVar{}
	tags.Set(str)

	if tags.String() != str {
		t.Fail()
	}
}
