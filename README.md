# Rgo
Connecting R and Go

This is a general repository for my attempts to use both R and Go in the same workflow for data analysis and whatnot. It contains 3 separate packages:
1. embarrasync:

    This is my attempt to make a Go package similar to R's parallel package (meaning it really only does embarrassingly parallel computation). It was a good way for me to get used to channels and Go's sync package, but it's also not helpful and you probably shouldn't use it. 
    
2. rfunc:

    This is simply a buch of convenient R functions written in Go, like mean, median, in, and unique. It's a handy library if you're used to R's functions but also want the type safety of Go. For that reason, rather than using something like interfaces to make the mean function work for all Go numeric types, the type is prepended to the function name, like Int64Mean.

 3. sexp:

    This is the most important package in the repo, as it allows one to connect R to Go, exporting Go functions as callable C functions. This makes it easier to use the good parts of Go (performance and non-embarrassing parallelization, for example) and the good parts of R (plotting and nice data analysis libraries, for example) into one workflow. It's the package I use and maintain the most. 
