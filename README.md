# Rgo
Connecting R and Go

This repository contains two packages:
1. `rfunc`:

    This is simply a bunch of convenient R functions written in Go, like mean, median, in, and unique. It's a handy library if you're used to R's functions but also want the type safety of Go. For that reason, rather than using something like interfaces to make the mean function work for all Go numeric types, the type is prepended to the function name, like `Int64Mean`.

    This would be very much improved if Go were to introduce generics.

 3. `sexp`:

    This is the most important package in the repo, as it allows one to connect R to Go, exporting Go functions as callable C functions. This makes it easier to use the good parts of Go (performance and non-embarrassing parallelization, for example) and the good parts of R (plotting and nice data analysis libraries, for example) into one workflow. 

This repository also contains an example of using `sexp` to export a library of simple Go functions that can be compiled and then called from R.