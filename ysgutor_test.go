package ysgutor

import (
	"testing"
	"fmt"
	"os"
)

func TestMain01(t *testing.T) {
	command := `echo "Hello World!"`
	fmt.Printf("Command: `%s`\tShell: nil\n", command)
	defer fmt.Println("==============================================\n\n")
	if y, err := New(command, nil); err == nil {
		if len(y.CommandArgs) == 1 {
			if y.CommandArgs[0] != `Hello World!` || y.CommandPath != `/bin/echo` {
				t.Error()
			}
		} else {
			fmt.Println(y.CommandArgs)
			t.Error()
		}
		y.StdOutWriter = os.Stdout
	} else {
		fmt.Println(err.Error())
		t.Error()
	}
	fmt.Println("OK")
}
