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

#Convert a vector of floats to a matrix from Go
Matricize <- function(InVec) {
	#Initialize the matrix
	m = matrix(nrow=(length(InVec)-1)/InVec[1], ncol=InVec[1])
	
	#Fill it in row by row
	colstart = 2
	colend = colstart + InVec[1] - 1
	for (r in 1:nrow(m)) {
		m[r,] = InVec[colstart:colend]
		colstart = colstart + InVec[1]
		colend = colend + InVec[1]
	}
	
	return(m)
}

#Convert a matrix to a vector for Go
Vectorize <- function(m) {
	#Make the vector
	v = vector(mode='double', length=length(m)+1)
	v[1] = ncol(m)
	
	ind = 2
	for (r in 1:nrow(m)) {
		for (c in 1:ncol(m)) {
			v[ind] = m[r,c]
			ind = ind + 1
		}
	}
	
	return(v)
}