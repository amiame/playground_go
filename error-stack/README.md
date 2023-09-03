# Objective
Figure out how simplify error messaging

## Motivation
Up until now, to make it easy tracing and error when it happens, I have been writing the function name, and the DDD layer
the error happens into the error message. And as the error gets sent up the layer, I would wrap that error with another
message indicating that function name and the layer it is in. This was a bit cumbersome.

After some Googling, I found out several Go packages to help remedy this situation:
- [runtime/debug](https://pkg.go.dev/runtime/debug@go1.21.0)
- [rotisserie/eris](https://pkg.go.dev/github.com/rotisserie/eris)
- [pkg/errors](https://pkg.go.dev/github.com/pkg/errors)

To put things in context, the Go package I have been using was the standard [errors](https://pkg.go.dev/errors) package.

## Conclusion
I find rotisserie/eris package the best, especially through the use of `eris.UnpackedError` to create custome message.
I believe that developers would appreciate not having to manually wrap every error messages everywhere they get returned to in the call stack.
Here's a comparison of the standard error code and eris's code:
### With standard error package
```go
// package1.go
func Package1Function1(arg1 int) error {
	if err := package1Function2(arg1); err != nil {
		return fmt.Errorf("this is an error in function 1: %w", err) <--- Here developer needs to remember to wrap
	}
	return nil
}

func package1Function2(arg int) error {
	if arg1 == 1 {
		return nil
	}
	return std_errors.New("this is an error in function 2")
}

// main.go (somewhere where we finally print the code)
// call Package1Function1
err = fmt.Errorf("this error is written in main: %w", err) <--- Also, developer needs to wrap here
fmt.Println(err)
for err != nil {
    fmt.Println(err)
    err = std_errors.Unwrap(err)
}
```
which produces:
```bash
this error is written in main: this is an error in function 1: this is an error in function 2
this is an error in function 1: this is an error in function 2
this is an error in function 2
```

### With rotisserie/eris package
```go
// package1.go
func Package1Function1(arg1 int) error {
	return package1Function2(arg1)  <--- No need to wrap
}

func package1Function2(arg int) error {
	if arg1 == 1 {
		return nil
	}
	return eris.New("this is an error in function 2")
}

// main.go
type Error struct {
	Message string   `json:"message"`
	Wraps   []string `json:"wrapping_messages,omitempty"`
	Calls   []string `json:"stack_calls,omitempty"`
}

func removeCurrentDirectory(path string) string {
	return strings.Replace(path, workingDir+"/", "", -1)
}

// call Package1Function1
unpackedErr := eris.Unpack(err)
customError := Error{Message: unpackedErr.ErrRoot.Msg}
if len(unpackedErr.ErrRoot.Stack) > 0 {
    firstStack := unpackedErr.ErrRoot.Stack[0]
    errorLocation := fmt.Sprintf("%s[L%d]: %s", removeCurrentDirectory(firstStack.File), firstStack.Line, firstStack.Name)
    customError.Message = errorLocation + ": " + unpackedErr.ErrRoot.Msg
}
for _, l := range unpackedErr.ErrChain {
    customError.Wraps = append(customError.Wraps, fmt.Sprintf("%s[L%d]: %s: %s", removeCurrentDirectory(l.Frame.File), l.Frame.Line, l.Frame.Name, l.Msg))
}
for _, s := range unpackedErr.ErrRoot.Stack {
    customError.Calls = append(customError.Calls, fmt.Sprintf("%s[L%d]: %s", removeCurrentDirectory(s.File), s.Line, s.Name))
}
u, _ = json.MarshalIndent(customError, "", "  ")
fmt.Printf("%v\n", string(u))
```
which produces
```bash
{
  "message": "package1/package1.go[L22]: package1.package1Function2: this is an error in function 2",
  "stack_calls": [
    "package1/package1.go[L22]: package1.package1Function2",
    "package1/package1.go[L12]: package1.Package1Function1",
    "main/main.go[L138]: main.main"
  ]
}
```

Granted the implementation for eris in the main code is longer, but what's more important is that it removes manual message composing needed to be done in
`Package1Function1`, and in main.go where the error does not originate from in the first place!
