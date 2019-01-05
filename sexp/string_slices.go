//This contains some functions to make passing vectors/slices of strings easier
package sexp

import(
	"strings"
)

//create a function that concatenates a slice of strings to a single string
func Slice2single(slice []string, delim string) (singlestring string) {
	singlestring = strings.Join(slice, delim)
	
	return singlestring
}

//this fuction separates a concatenated string into a slice
func String2slice(singlestring, delim string) (slice []string) {
	slice = strings.Split(singlestring, delim)
	
	return slice
}