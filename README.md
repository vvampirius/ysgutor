# ysgutor

   Just simple wrapper for os/exec with commandline parsing and callbacks.

## Examples:

```
package main
import (
	"github.com/vvampirius/ysgutor"
	"os"
)
func main() {
	y, _ := ysgutor.New(`echo "Hello World!"`, nil)
	y.StdOutWriter = os.Stdout
	y.Execute()
}
```

```
package main
import (
	"github.com/vvampirius/ysgutor"
	"os"
	"context"
	"time"
	"fmt"
)
func main() {
	y, _ := ysgutor.New(`echo \"Hello World!\" ; sleep 10`, `/bin/sh -c "%s"`)
	y.OnPreStart = func(yy *ysgutor.Ysgutor) bool {
		yy.Cmd.Stdout = os.Stdout
		return true
	}
	y.Context, _ = context.WithTimeout(context.Background(), 2*time.Second)
	y.OnExit = func(yy *ysgutor.Ysgutor, _ error) {
		fmt.Println(yy.ExitedAt.Sub(yy.StartedAt))
	}
	y.Execute()
}
```