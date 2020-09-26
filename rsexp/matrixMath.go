package rsexp

// dotProduct calculates the dot product of two vectors. It is used to make matrix multiplication easier.
func dotProduct(a, b []float64) (f float64, err error) {
	if len(a) != len(b) {
		return f, SizeMismatch
	}

	for i, av := range a {
		f += av * b[i]
	}
	return f, nil
}

// MatrixMultiply performs a matrix multiplication of two matrices. This is not an element-wise multiplication, but
// a true multiplication as defined in elementary linear algebra. In matrix multiplication, order
// matters. Two matrices A and B can only be multiplied if A has the same number of rows as B has number of columns. If
// the dimensions of the input matrices do not allow for a multiplication, a SizeMismatch error is returned.
func MatrixMultiply(A, B *Matrix) (C *Matrix, err error) {
	// checks to ensure matrix quality
	if !(A.isSizeValid() && B.isSizeValid()) {
		return nil, ImpossibleMatrix
	}
	if A.Ncol != B.Nrow {
		return nil, SizeMismatch
	}

	// create the matrix with the right dimensions
	C, err = CreateZeros(A.Nrow, B.Ncol)
	if err != nil {
		return C, err
	}

	for i := 0; i < A.Nrow; i++ { // this indexes the row
		for j := 0; j < B.Ncol; j++ { // this indexes the column
			aVec, err := A.GetRow(i)
			if err != nil {
				return C, err
			}
			bVec, err := B.GetCol(j)
			if err != nil {
				return C, err
			}
			dp, err := dotProduct(aVec, bVec)
			if err != nil {
				return C, err
			}
			err = C.SetInd(i, j, dp)
			if err != nil {
				return C, err
			}
		}
	}

	return C, nil
}

// MatrixAdd adds two matrices. Matrix addition is done by adding each element of the two matrices together, so
// they must be of identical size. If they are not, a SizeMismatch error will be returned.
func MatrixAdd(A, B *Matrix) (C *Matrix, err error) {
	// checks to ensure matrix quality
	if !(A.isSizeValid() && B.isSizeValid()) {
		return nil, ImpossibleMatrix
	}
	if !(A.Nrow == B.Nrow && A.Ncol == B.Ncol) {
		return nil, SizeMismatch
	}

	C = &Matrix{Nrow: A.Nrow, Ncol: B.Ncol, Data: make([]float64, len(A.Data))}

	for i, av := range A.Data {
		C.Data[i] = av + B.Data[i]
	}

	return C, nil
}

// AddConstant adds a constant to every element of a matrix. There is no SubtractConstant method. To subtract a
// constant N from a matrix, add its negative, -N.
func (m *Matrix) AddConstant(c float64) {
	for i, _ := range m.Data {
		m.Data[i] += c
	}
}

// MultiplyConstant multiplies each element of a matrix by a constant. There is no DivideConstant method. To divide a
// matrix by a constant N, multiply it by its reciprocal, 1/N.
func (m *Matrix) MultiplyConstant(c float64) {
	for i, _ := range m.Data {
		m.Data[i] *= c
	}
}
