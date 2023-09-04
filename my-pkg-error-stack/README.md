# `errors` package

## Background
This package was made to improve error tracing.
The problem was that for the sake of making it easier to trace bugs, we wrote information about:
- the function name the error occurred,
- the package the error occured

The reason was to assist developers in tracing the source of the bug.

This practice however has several drawbacks:
- you have to write error messages for every function the error was returned to
- you have to write the error messages manually

This made the practice prone to human error.

## This package's solution
By utilizing the standard `runtime` package, this package retrieves the call stack from the point the error occured.
For instance:
```go
package main

import (
	"fmt"
	"myPkg/errors"
)

func funcA() error {
	return errors.New("hey you!")
}

func main() {
	err := funcA()
    upErr := errors.Unpack(err)
	u, _ := json.MarshalIndent(upErr, "", "  ")
	fmt.Printf("%v\n", string(u))
}

/* output:
{
  "calls": [
    {
      "message": "hey you!",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.funcA",
      "line": 10
    },
    {
      "message": "",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.main",
      "line": 14
    }
  ]
}
*/
```

Moreover, if in that call stack, the developer called an `errors.Wrap`, that would be printed out as well.
For example:
```go
package main

import (
	"encoding/json"
	"fmt"
	"myPkg/errors"
)

func funcA() error {
	return errors.New("hey you!")
}

func main() {
	err := funcA()
	errors.Wrap(err, "This is wrap")
	errors.Wrapf(err, "This is another %s", "wrap")
	upErr := errors.Unpack(err)
	u, _ := json.MarshalIndent(upErr, "", "  ")
	fmt.Printf("%v\n", string(u))
}
/* output:
{
  "calls": [
    {
      "message": "hey you!",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.funcA",
      "line": 10
    },
    {
      "message": "",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.main",
      "line": 14
    },
    {
      "message": "This is wrap",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.main",
      "line": 15
    },
    {
      "message": "This is another wrap",
      "file": "/Users/amir_c/work/playground_go/worktrees/others/my-error-stack-package/my-pkg-error-stack/example_usecase/main/main.go",
      "function": "main.main",
      "line": 16
    }
  ]
}
*/
```

## Other conveniences

### Composing error logs in your own style
By utilizing the `UnpackedError` struct, the developer can compose their own logging structure.
For example, they can shorten the `file` value by cut out the project directory from the path:
```go
package main

import (
	"fmt"
	"myPkg/errors"
	"os"
	"strings"
)

var workingDir string

func init() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
	workingDir = dir
}

func funcA() error {
	return errors.New("hey you!")
}

func removeCurrentDirectory(path string) string {
	return strings.Replace(path, workingDir+"/", "", -1)
}

func main() {
	err := funcA()
	errors.Wrap(err, "I'm wrapped")
	errors.Wrap(err, "I'm wrapped2")

	upErr := errors.Unpack(err)
	for i, call := range upErr.Calls {
		if call.Msg != "" {
			fmt.Printf("%d) %s[%d]: %s: %s\n", i, removeCurrentDirectory(call.File), call.Line, call.Function, call.Msg)
		} else {
			fmt.Printf("%d) %s[%d]: %s\n", i, removeCurrentDirectory(call.File), call.Line, call.Function)
		}
	}
}
/* output:
0) example_usecase/main/main.go[22]: main.funcA: hey you!
1) example_usecase/main/main.go[30]: main.main
2) example_usecase/main/main.go[31]: main.main: I'm wrapped
3) example_usecase/main/main.go[32]: main.main: I'm wrapped2
*/
```
