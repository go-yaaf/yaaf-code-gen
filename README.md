# yaaf-code-gen
Code generator utility

This library is used to generate documentation / code / script files based on the go source comments in your Go project.
This library is based on `go/parser` and `go/ast` packages combine with template engine to parse go source files
and based on special tokens in the comments, to generate text files for various cases.

Common use cases are:
* Generate API documentation
* Generate API client library for various programing languages
