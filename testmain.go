// +ignore
/*
<!--
Copyright (c) 2020 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Cannot import main: a Go Module gotcha"
description = "How using go test on a main package can clash with module naming"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2020-11-28"
draft = "true"
categories = ["Go Ecosystem"]
tags = ["test", "modules", ""]
articletypes = ["Background"]
+++

Two questions for you: Do you name an app module simply "main"? And do you happen to write tests for a main package? If so, you are in big trouble! (Ouch, that's quite clickbait-ey, isn't it?) Well, the world is not exactly going to end; however, you might encounter an unexpected error that is hard to track down.

<!--more-->


## Ingredient #1: a main module

Go Modules apply to library packages as well as to executable packages. Imagine that you create a new executable project. Usually you likely start by running

```sh
go mod init main
```

in the root directory of your project. Since this is going to become an app and not a library, you choose the name "main" as the module "path". After all, nothing and nobody would want to import a main package.

With this in mind, let's move on to...


## Ingredient #2: a test file

So your main file has grown large but not large enough to move some code out of package `main` into library packages. Still, a couple of rather complex functions have already accumulated there, and so you decide to create a unit test for package `main`.

After all, who says that unit tests are restricted to library packages? And this worked before (read: back in the times we all still used `GOPATH`), so let's go ahead and write some tests in `main_test.go`.

An example:
*/
// The go.mod file
module main

go 1.15
/*

*/
// This is the package to test, file "main.go"
package main


func main() {
	square(4)
}

func square(n int) int {
	return n * n
}

/*
*/
// And this is our test file, "test_main.go"
package main

import "testing"

func Test_square(t *testing.T) {
	t.Log("Test successful")
}
/*

## Now things start to turn weird

Finally, with all tests in place, you run

```sh
> go test
```

inside the project directory, but instead of the usual output, you get this:

```
# main.test
/var/folders/_m/dgnkqt8d3j10svk5c06px4vc0000gn/T/go-build306511963/b001/_testmain.go:13:2: cannot import "main"
FAIL    main [build failed]
FAIL
```

Your test file uses package `main`, and so there is no `import main` anywhere in your test file. After all, no one would ever want to import package `main`, right?

Well, `go test` apparently tries exactly that. Or at least it seems so.

However, if you run instead

```sh
> go test *.go
```

all is good, and you get:

```
ok      command-line-arguments  0.293s
```

So the problem only occurs for a plain

`go test`

and for

`go test .`

What is going on here?


## The Web knows it all - or does it?


Even a thorough search reveals only vague information from some obscure corners of the Go repo on GitHub:

https://github.com/golang/go/issues/10738#issuecomment-99732939

> The problem is that the test driver (_testmain.go) needs to import the package and if the import path is "main", gc will refuse to import it.
>
> The reason is simple, the packages in Go program must have unique import paths, but as the main package must have the import path "main", you can't import another "main" package.


https://github.com/golang/go/issues/28514#issuecomment-457734995

> go build works with the module named main because nothing ends up importing the resulting package.
>
> go test does not work because the generated test package needs to import the package under test in order to actually run the tests.


At least, this is a staring point. None of the above answers, however, can explain the dichotomy of `go test` versus `go test *.go`.


## Trying to get the full picture

Ok, let's put the pieces together.

### The test driver

Go's test driver needs to import the package it is going to test. This even applies if the test code is in the same package as the code that it tests (also known as white-box test). Fair enough. Still, how come the driver even *wants* to import `main` as a package? Especially when the test functions are inside package `main` already.


### Go modules and the import path

Another comment from issue #28514 reveals that the module declaration in `go.mod` influences the behavior of `go test`: 

> The error went away when I replaced the package name (main) in go.mod with the module path (my/project).

Let's try that. Or even simpler, let's just *rename* the module, rather than adding a path.

*/
// The go.mod file
module mymain

go 1.15
/*

Aaand (drum rolls)...

```sh
> go test
PASS
ok      mymain  0.355s
```

Now the test of package `main` works again. Note that the package name was not changed, only the module name. 


## Conclusion

Now we have two options for addressing the problem.

1. Do not write tests for package main. Main is for integrating your code units (a.k.a packages) and usually would not contain unit-testable code. IF you have test-worthy code in main (for unit tests, that is; not for integration tests), consider moving it into a library package.

2. Or rename your "main" module. You can choose a name of your liking, and you do not even need to add an import path prefix, as package `main` is not supposed to be fetched by an `import` directive. So for your next project, try something like `go mod init jumpingjehoshaphat`.


**Happy coding!**

*/
