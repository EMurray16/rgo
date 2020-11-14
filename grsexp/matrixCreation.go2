package rsexp

// CreateZeros creates a matrix of the given dimensions in which every element is 0. If the given dimensions
// are nonsensical (negative, for example) it will return an InvalidIndex error.
func CreateZeros(Nrow, Ncol int) (*Matrix, error) {
	if Nrow <= 0 || Ncol <= 0 {
		return nil, InvalidIndex
	}
	//create the vector
	v := make([]float64, Nrow*Ncol)
	//create the matrix
	m := Matrix{Nrow: Nrow, Ncol: Ncol, Data: v}

	return &m, nil
}

// CreateIdentity creates an identity matrix, which is always square by definition, of the input dimension.
// An identity matrix is a matrix will all 0s, except for having a 1 in each element of the diagonal. If the
// given size is impossible, it will return an InvalidIndex error.
func CreateIdentity(size int) (*Matrix, error) {
	// start by creating a zero matrix
	outMat, err := CreateZeros(size, size)
	if err != nil {
		return outMat, err
	}
	// now fill in the 1s
	for i := 0; i < size; i++ {
		err := outMat.SetInd(i, i, 1)
		if err != nil {
			return outMat, err
		}
	}
	return outMat, nil
}

// NewMatrix creates a new matrix given a vector of data. The number of rows and columns must be provided, and it
// assumes the data is already in the order a Matrix should be, with column indexes adjacent. In other
// words, the data vector should be a concatenation of several vectors, one for each column. NewMatrix makes a copy
// of the input slice, so that changing the slice later will not affect the data in the matrix. If the provided
// dimensions don't match the length of the provided data, an ImpossibleMatrix error will be returned.
func NewMatrix(Nrow, Ncol int, data []float64) (*Matrix, error) {
	if Nrow < 0 || Ncol < 0 {
		return nil, InvalidIndex
	}
	if len(data) != Nrow*Ncol {
		return nil, ImpossibleMatrix
	}

	outMat := &Matrix{Nrow: Nrow, Ncol: Ncol}
	outMat.Data = make([]float64, Nrow*Ncol)
	copy(outMat.Data, data)
	return outMat, nil
}

// CopyMatrix creates an exact copy of an existing matrix. The copies are independent, so that the output matrix can
// be changed without changing the input matrix and vice versa.
func CopyMatrix(in Matrix) (out Matrix) {
	out.Nrow = in.Nrow
	out.Ncol = in.Ncol
	out.Data = make([]float64, len(in.Data))
	for i, f := range in.Data {
		out.Data[i] = f
	}
	return out
}

// CreateTranspose creates a new matrix which is a transpose of the input matrix. The output matrix is created from
// a copy of the input matrix such that they can be altered independently.
func (m *Matrix) CreateTranspose() *Matrix {
	// create the matrix with the new dimensions
	mt := Matrix{Nrow: m.Ncol, Ncol: m.Nrow, Data: make([]float64, len(m.Data))}

	// each row of the new matrix is a column of the old one
	for i := 0; i < m.Ncol; i++ {
		rowToSet, _ := m.GetCol(i)
		mt.SetRow(i, rowToSet)
	}
	return &mt
}
