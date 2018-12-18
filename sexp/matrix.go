//This adds support for matrices and vectors of strings
package sexp

import "fmt"

//create a type for a matrix of floats
type Matrix struct {
	//indexed as [row][col]
	Mat [][]float64
}

//We could make a matrix of ints here, but I don't think it'd ever be useful *shrug*

//create a couple methods for convenience (similar to R)
func (M Matrix) Length() (l int) {
	for _, col := range M.Mat {
		l += len(col)
	}
	return l
}
func (M Matrix) Nrow() (n int) {
	n = len(M.Mat)
	return n
}
func (M Matrix) Ncol() (n int) {
	n = len(M.Mat[0])
	return n
}

//create a function that initializes a matrix
func CreateMatrix(nrow, ncol int) (M Matrix) {
	M.Mat = make([][]float64, nrow)
	for r, _ := range M.Mat {
		M.Mat[r] = make([]float64, ncol)
	}
	return M
}

//create a method that "vectorizes" a matrix
func (M Matrix) Vectorize() (v []float64) {
	//start by making the vector
	v = make([]float64, M.Length()+1)
	//the first vector value is the length of a row (or number of columns)
	v[0] = float64(M.Ncol())

	//now fill in the rest of the vector
	ind := 1
	for _, row := range M.Mat {
		for _, f := range row {
			v[ind] = f
			ind++
		}
	}

	return v
}

//create a method that unvectorizes a vectorized matrix
func Matricize(v []float64) (M Matrix) {
	ncol := int(v[0])
	nrow := (len(v) - 1) / ncol

	//make the matrix
	M = CreateMatrix(nrow, ncol)

	//fill in the matrix element by element
	for ri, r := range M.Mat {
		for ci, _ := range r {
			ind := 1 + (ri * ncol) + ci
			M.Mat[ri][ci] = v[ind]
		}
	}

	return M
}

//create a function that reads a vectorized matrix into a Go matrix
func AsMatrix(s GoSEXP) (m Matrix) {
	vec := AsFloats(s)
	fmt.Println(vec)
	m = Matricize(vec)
	return m
}

//create a function that function that converts a matrix to a sexp
func Mat2sexp(M Matrix) (s GoSEXP) {
	//vectorize the matrix
	vec := M.Vectorize()
	//now write the vector
	s = Float2sexp(vec)
	return s
}
