#This contains helper functions for going between R and Go

# ParseGoMatrix is a convenience function for creating an R matrix from a list returned from Go that was created by
# the Matrix2sexp function.
rsexp.ParseGoMatrix <- function(rsexpOutput) {
    # The matrix is the second element of the list, while the first is an integer vector c(nrow,ncol)
    outMat = matrix(data=rsexpOutput[[2]], nrow=rsexpOutput[[1]][1], ncol=rsexpOutput[[1]][2])
    return(outMat)
}