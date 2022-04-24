# snippetbox

### Notes:
- `go.mod` - project is a module (unique)
- Running equivalence
```
go run . 
go run github.com/kohrongying/snippetbox
go run main.go
```
- Go's servemux supports
1. fixed path (eg. /snippet/view, /snippet/create)
2. subtree path (end with trailing slash) (eg. home)