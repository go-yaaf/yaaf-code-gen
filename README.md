# yaaf-code-gen

[![Build](https://github.com/go-yaaf/yaaf-code-gen/actions/workflows/build.yml/badge.svg)](https://github.com/go-yaaf/yaaf-code-gen/actions/workflows/build.yml)

**yaaf-code-gen** is an artifacts generation package.

## Core Philosophy
During the process of product development, there are several artifacts that their content is derived from 
the code base, some common use cases are:
- Generate documentation
- Generate tests
- Generate client-side libraries for REST APIs for the consumers (e.g. TypeScript library, Javascript etc)


To make this artifacts always aligned with the latest version and avoid manual updates, the purpose of this
package is to generate these artifacts based on pre-defined annotations in the code (similar to Java DocGen concepts)

## Implementation
This library is based on `go/parser` and `go/ast` packages combine with template engine to parse go source files.
The process includes the following steps:
1. Parsing the source code of the provided paths and create a meta model describing the building blocks of the application: data structures, enums, services ...
2. Processing the meta model and using templates to generate artifacts (files) based on this model


## Installation

To add `yaaf-code-gen` to your project, use `go get`:
```bash
go get -u github.com/go-yaaf/yaaf-code-gen
```

## Developer Guide

This guide provides an overview of the core components and how to use them.



