#This contains helper functions for going between R and Go

#Convert an R string to Go bytes
GoEncode <- function(InString) {
	InBytes = as.integer(charToRaw(InString))
	return(InBytes)
}

#Convert Go bytes to an R string
GoDecode <- function(InBytes) {
	InString = rawToChar(as.raw(InBytes))
	return(InString)
}