# snippetbox

### Notes:
- `go.mod` - project is a module (unique)
- Running equivalence
```
go run . 
go run github.com/kohrongying/snippetbox
go run main.go
```
#### Go's servemux 
- supports
    1. fixed path (eg. /snippet/view, /snippet/create)
    2. subtree path (end with trailing slash) (eg. home)
- Servemux dispatch request to the handler with longest corresponding pattern. --> can register patterns and handlers in any order
- Request URL auto sanitized eg. /foo/bar/./../baz --> /foo/bar/baz 301
- Subtree path request without / eg. /foo --> 301 to /foo/
- Pattern matching --> host specific checked first, then non-host specific patterns next.

#### http.ResponseWriter
- if do not explicitly state w.WriteHeader (eg. w.WriteHeader(405)), error will return with 200 
```go
w.WriteHeader(405)
w.Write([]byte("Method not supported"))

#same as
http.Error(w, "Method not supported" , 405)
```
- Go sends system generated response headers (Date,  Content-Length and Content-Type).
Headers can be manipulated by
`w.Headers().Get(<key>)`, `w.Headers().Set(<key>, <value>)`, `w.Headers().Add(<key>, <value>)`, `w.Headers().Del(<key>)`
If suppress the sys generataed ones: `w.Header()["Date"] = nil`
- fmt.Fprintf takes a io.Writer but able to pass http.ResponseWRiter object, as io.Writer is an interface and http.ResponseWriter satisfies as it as a w.Write method. 

#### Requests
All incoming HTTP requests are served in their own goroutine --> concurrency needs to be handled. Race conditions.

#### Structure
cmd - applicable specific code
| web - only one executable
    | handlers.go, main.go

internal - non-applicable specific code

ui/html, ui/static

#### Logging
- Recommened to use Panic and Fatal in main() and not elsewhere
- Custom loggers are concurrency safe. Share a single logger across multiple goroutines
- Log output to standard streams and redirect output to file at runtime eg. `go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log`