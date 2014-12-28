Generate Functional Options
=======

Functional options can be used to create elegant APIs in Go.

This project will install the command `optioner` which is a tool to generate functional options. The `optioner` command
is intended to be used with `go generate`.

The idea of functional options was first introduced in this blog post: [Self Referential Functions and Design] by Rob Pike(http://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html). More recently, Dave Cheney presented [Functional Options for Friendly APIs](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

How ot use:

```
go get github.com/akualab/optioner
go intall github.com/akualab/optioner
cd $GOPATH/github.com/akualab/optioner/example
go generate
go test
```

Documentation can be found at http://godoc.org/github.com/akualab/optioner
