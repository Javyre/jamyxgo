# Jamyxgo
Go client library for jamyxer.
This library was created for development of [jamyxui](github.com/javyre/jamyxui).
This is a wrapper of some [jamyxer](github.com/javyre/jamyxer) commands passed to the server via tcp.

# Installation
    go get github.com/javyre/jamyxgo

# Usage
```golang
import "github.com/javyre/jamyxer"

func main() {
    session := jamyxgo.Session{}

    // default port for jamyxer is 2909
    session.Connect("127.0.0.1", 2909)

    // Do stuff...
}
```

For documentation check out http://godoc.org/github.com/Javyre/jamyxgo
