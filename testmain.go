/*
<!--
Copyright (c) 2019 Christoph Berger. Some rights reserved.

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
description = "Why using go test on a main package is a bad idea"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2020-10-01"
draft = "true"
categories = ["Go Ecosystem"]
tags = ["test", "modules", ""]
articletypes = ["Background"]
+++

Go unit tests do work with a main package -- until you use Go Modules. Read what happens then, and why Go modules are not the culprit.

<!--more-->

## A rare encounter

The following scenario is supposed to happen only rarely, and only to Go newcomers. However, the time needed to find the cause of the problem can be significant, hence I decided to describe and explain the problem here in case anyone runs into it


## Ingredient #1: a main module

Go Modules apply to library packages as well as to executable packages. So when you create a new executable project, you will likely start with running

```sh
go mod init main
```

in the root directory of your project. Using "main" here as the module name instead of a global import path like "github.com/appliedgo/mainproject" makes totally sense, as a main package cannot be imported.

With this in mind, let's move on to


## Ingredient #2: a test file

So your main file has grown large but not large enough to move some code out of package `main` into library packages. Still, there are a couple of rather complex functions inside, so you decide to create a unit test for package `main`.

After all, who says that unit tests are restricted to library packages? And this worked before, so let's go ahead and write some tests in `main_test.go`.


## Now things turn weird

Finally, with all tests in place, you run

```sh
go test .
```

but instead of the usual output, you get this:

```
# main.test
/var/folders/_m/dgnkqt8d3j10svk5c06px4vc0000gn/T/go-build306511963/b001/_testmain.go:13:2: cannot import "main"
FAIL    main [build failed]
FAIL
```

Your test file uses package `main`, and so there is no `import main` anywhere in your test file. After all, no one would ever want to import package `main`, right?

Well, `go test` apparently does. Or at least it seems so. What is going on here?


##


* ON the internet, only some low level explanations in obscure locations

https://github.com/golang/go/issues/10738#issuecomment-99732939

> The problem is that the test driver (_testmain.go) needs to import the package and if the import path is "main", gc will refuse to import it.
>
> The reason is simple, the packages in Go program must have unique import paths, but as the main package must have the import path "main", you can't import another "main" package.


https://github.com/golang/go/issues/28514#issuecomment-457734995

> go build works with the module named main because nothing ends up importing the resulting package.
>
> go test does not work because the generated test package needs to import the package under test in order to actually run the tests.


## Conclusion

* Do not write tests for package main. Main is for integrating your code units (a.k.a packages) and should not contain unit-testable code.



## The code
*/

// ## Imports and globals
package main

/*
## How to get and run the code

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

    go get -d github.com/appliedgo/TODO:

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/TODO:

Step 3. Run the binary.

    go run TODO:.go


## Odds and ends
## Some remarks
## Tips
## Links


**Happy coding!**

*/
