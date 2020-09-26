package rsexp

// AreMatricesEqual returns true if the input matrices are of the same dimension and have identical data vectors. It's
// important to note that this function uses strict equality - even if elements of two matrices differ by floating
// point error, it will return false.
func AreMatricesEqual(A, B Matrix) bool {
	if A.Ncol != B.Ncol {
		return false
	}
	if A.Nrow != B.Nrow {
		return false
	}

	// This check should be superfluous if Nrow and Ncol are the same, but because
	// the package exports a struct it's possible for the data to become misaligned with
	// the matrix header information
	if len(A.Data) != len(B.Data) {
		return false
	}

	// TODO: I thought there was a better way to check for equality than this?
	for i, v := range A.Data {
		if v != B.Data[i] {
			return false
		}
	}

	return true
}

// AreMatricesEqualTol is the same as AreMatricesEqual, except the data vectors are checked in relation to the
// input tolerance allowed. If any elements differ by more than the tolerance, this function will return false.
func AreMatricesEqualTol(A, B Matrix, tolerance float64) bool {
	if A.Ncol != B.Ncol {
		return false
	}
	if A.Nrow != B.Nrow {
		return false
	}

	// This check should be superfluous if Nrow and Ncol are the same, but because
	// the package exports a struct it's possible for the data to become misaligned with
	// the matrix header information
	if len(A.Data) != len(B.Data) {
		return false
	}

	// TODO: I thought there was a better way to check for equality than this?
	for i, v := range A.Data {
		if v-B.Data[i] > tolerance || B.Data[i]-v > tolerance {
			return false
		}
	}

	return true
}

// GetRow gets the row of the matrix specified by the provided index, using 0-based indexing. The first row of a matrix
// is index 0, even though it may be more intuitive that it should be 1. If the input index is too big, it will return
// a IndexOutOfBounds error. If you get this error, there's a good chance it's just an off-by-one error. The resulting
// slice does not point to the matrix itself, so it can be edited without altering the matrix.
func (m *Matrix) GetRow(ind int) ([]float64, error) {
	if !m.isSizeValid() {
		return nil, ImpossibleMatrix
	}
	if ind < 0 {
		return nil, InvalidIndex
	}
	if ind >= m.Nrow {
		return nil, IndexOutOfBounds
	}

	f := make([]float64, m.Ncol)
	//the row are not adjacent - row 1 is element 0, Ncol, 2*Ncol, etc.
	for offsetCount := 0; offsetCount < m.Ncol; offsetCount++ {
		f[offsetCount] = m.Data[offsetCount*m.Nrow+ind]
	}
	return f, nil
}

// SetRow sets the row of the matrix, specified by the input index, to match the data provided. If the provided data
// is not of the same length as the number of columns in the matrix, it will return a SizeMismatch error.
func (m *Matrix) SetRow(ind int, data []float64) error {
	if !m.isSizeValid() {
		return ImpossibleMatrix
	}
	if len(data) != m.Ncol {
		return SizeMismatch
	}
	if ind < 0 {
		return InvalidIndex
	}
	if ind >= m.Nrow {
		return IndexOutOfBounds
	}

	for offsetCount := 0; offsetCount < m.Ncol; offsetCount++ {
		m.Data[offsetCount*m.Nrow+ind] = data[offsetCount]
	}

	return nil
}

// AppendRow appends a row onto an existing matrix and updates the dimension metadata accordingly. If the length of the
// provided row is not equal to the number of columns in the matrix, it will return a SizeMismatch error.
func (m *Matrix) AppendRow(data []float64) error {
	if !m.isSizeValid() {
		return ImpossibleMatrix
	}
	if len(data) != m.Ncol {
		return SizeMismatch
	}

	//make a dummy slice of 0s for each added element
	dummy := make([]float64, m.Ncol)
	m.Data = append(m.Data, dummy...)

	//now we can safely index the Ncol field without lying
	m.Nrow++

	//loop through the end of each row and insert the data
	var sliceind int = 0 //tracks data index
	for i := m.Nrow - 1; i < len(m.Data); i += m.Nrow {
		copy(m.Data[i+1:], m.Data[i:]) //this opens the blank space in the slice
		m.Data[i] = data[sliceind]
		sliceind++
	}

	return nil
}

// GetCol gets the column of the matrix specified by the provided index, using 0-based indexing. The first column of a
// matrix is index 0, even though it may be more intuitive that it should be 1. If the input index is too big, it will
// return a IndexOutOfBounds error. If you get this error, there's a good chance it's just an off-by-one error. The
// resulting slice does not point to the matrix itself, so it can be edited without altering the matrix.
func (m *Matrix) GetCol(ind int) ([]float64, error) {
	if !m.isSizeValid() {
		return nil, ImpossibleMatrix
	}
	if ind < 0 {
		return nil, InvalidIndex
	}
	if ind >= m.Ncol {
		return nil, IndexOutOfBounds
	}

	//the column indices are adjacent, so we just need to find the indexes
	f := make([]float64, m.Nrow)

	//now loop through the matrix and add elements one by one
	startInd := ind * m.Nrow
	endInd := (ind + 1) * m.Nrow
	copy(f, m.Data[startInd:endInd])
	return f, nil
}

// SetCol sets the column of the matrix, specified by the input index, to match the data provided. If the length of the
// provided column is not of the same as the number of rows in the matrix, it will return a SizeMismatch error.
func (m *Matrix) SetCol(ind int, data []float64) error {
	if !m.isSizeValid() {
		return ImpossibleMatrix
	}
	if len(data) != m.Nrow {
		return SizeMismatch
	}
	if ind < 0 {
		return InvalidIndex
	}
	if ind >= m.Ncol {
		return IndexOutOfBounds
	}

	for i := 0; i < m.Nrow; i++ {
		m.Data[ind*m.Nrow+i] = data[i]
	}

	return nil
}

// AppendCol appends a column onto an existing matrix and updates the dimension metadata accordingly. If the provided
// data column is not equal to the number of rows in the matrix, it will return a SizeMismatch error.
func (m *Matrix) AppendCol(data []float64) error {
	if !m.isSizeValid() {
		return ImpossibleMatrix
	}
	if len(data) != m.Nrow {
		return SizeMismatch
	}

	// because column indices are adjacent, we just need to use append
	m.Data = append(m.Data, data...)
	m.Ncol++
	return nil
}

// This method returns the value in the element of the matrix defined by the inputs.
func (m *Matrix) GetInd(row, col int) (float64, error) {
	if !m.isSizeValid() {
		return 0, ImpossibleMatrix
	}

	var output float64
	i := col*m.Nrow + row
	if i >= len(m.Data) {
		return 0, IndexOutOfBounds
	}
	if i < 0 || row < 0 || col < 0 {
		return 0, InvalidIndex
	}

	output = m.Data[i] // make sure we send back a real copy, not a pointer

	return output, nil
}

// This method alters the value in the element of the matrix defined by the inputs to the input value.
func (m *Matrix) SetInd(row, col int, data float64) error {
	if !m.isSizeValid() {
		return ImpossibleMatrix
	}

	i := col*m.Nrow + row
	if i >= len(m.Data) {
		return IndexOutOfBounds
	}
	if i < 0 || row < 0 || col < 0 {
		return InvalidIndex
	}

	m.Data[i] = data
	return nil
}
