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

### psql
```sql
create database snippetbox
\c snippetbox
CREATE TABLE snippets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    expires TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    now(),
    NOW() + INTERVAL '5 DAY'
),
 (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    now(),
    NOW() + INTERVAL '5 DAY'
),
(
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    now(),
    NOW() + INTERVAL '5 DAY'
);
CREATE USER web;
grant connect on database snippetbox to web;
grant select, insert, update, delete on all tables in schema public to web;
```

#### Get psql driver
- Command: `go get github.com/lib/pq@v1.10.5`
- Updates go.mod and creates go.sum files
- Import path for driver prefixed with underscore. our `main.go` file does not use anything in that package, without the _, Go compiler will raise error. We need driver's init() function to run so that it can register itself with database/sql package.
