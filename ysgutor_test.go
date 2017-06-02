package ysgutor

import (
	"testing"
)

func TestMain(t *testing.T) {
	if y, err := New(`echo "Hello World!"`, nil); err == nil {
		//fmt.Println(y.CommandArgs)
		if len(y.CommandArgs) == 1 {
			if y.CommandArgs[0] != `Hello World!` || y.CommandPath != `/bin/echo` {
				t.Error()
			}
		} else {
			t.Error()
		}
	}
}
