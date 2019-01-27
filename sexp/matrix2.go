//This is a matrix implementation for the sexp package
package sexp

// import "fmt" //used for troubleshooting from time to time

//create the matrix type
type Matrix struct {
	Nrow, Ncol int
	V []float64
}

//this function creates a zero-value matrix
func CreateZeros(nrow, ncol int) *Matrix {
	//create the vector
	v := make([]float64, nrow*ncol)
	//create the matrix
	m := Matrix{nrow, ncol, v}
	
	return &m
}

//this method returns an entire row of the matrix as a slice
func (m *Matrix) GetRow(ind int) (f []float64) {
	//the row indexes are next to each other
	startind := (ind - 1) * m.Ncol
	f = m.V[startind:(startind+m.Ncol)]
	return f
}

//this method returns an entire column of the matrix as a slice
func (m *Matrix) GetCol(ind int) (f []float64) {
	//the col indexes are not contiguous, so preallocate the slice
	f = make([]float64, m.Nrow)
	
	//now loop through the matrix and add elements one by one
	sliceind := 0
	for i := ind-1; sliceind < m.Nrow; i += m.Ncol{
		f[sliceind] = m.V[i]
		sliceind++
	}
	
	return f
}

//this method gets a single index of the matrix
func (m *Matrix) GetInd(rowi, coli int) float64 {
	i := (rowi - 1) * m.Ncol + (coli - 1)
	return m.V[i]
}

//this method sets the entire row of a matrix
func (m *Matrix) SetRow(ind int, data []float64) () {
	startind := (ind - 1) * m.Ncol
	//edit the row in a loop
	dataind := 0
	for i := startind; i < startind+m.Ncol; i++ {
		m.V[i] = data[dataind]
		dataind++
	}
}

//this method sets the entire column of a matrix
func (m *Matrix) SetCol(ind int, data []float64) () {
	sliceind := 0
	for i := ind-1; sliceind < m.Nrow; i += m.Ncol {
		m.V[i] = data[sliceind]
		sliceind++
	}
}

//this method sets a single index of a matrix
func (m *Matrix) SetInd(rowi, coli int, data float64) () {
	i := (rowi - 1) * m.Ncol + (coli - 1)
	m.V[i] = data
}

//this method adds to all elements of a row
func (m *Matrix) AddRow(ind int, data []float64) () {
	//start by getting the row of the matrix
	row := m.GetRow(ind)
	
	//now add each element of data to row
	for i, v := range row {
		row[i] = v + data[i]
	}
	
	//now set that row
	m.SetRow(ind, row)
}

//this method adds to all elements of a column
func (m *Matrix) AddCol(ind int, data []float64) () {
	//start by getting the col of the matrix
	col := m.GetCol(ind)
	
	//now add each element of data to col
	for i, v := range col {
		col[i] = v + data[i]
	}
	
	//now set that col
	m.SetCol(ind, col)
}

//this method adds to a single index of the matrix
func (m *Matrix) AddInd(rowi, coli int, data float64) () {
	//get the index
	i := (rowi - 1) * m.Ncol + (coli - 1)
	m.V[i] = m.V[i] + data
}

//this method appends a row onto a matrix
func (m *Matrix) AppendRow(data []float64) () {
	//this is easy
	m.Nrow++
	m.V = append(m.V, data...)
}

//this method appends a col onto a matrix
func (m *Matrix) AppendCol(data []float64) () {
	//make a dummy slice of 0s for each added element
	dummy := make([]float64, m.Nrow)
	m.V = append(m.V, dummy...)
	
	//now we can safely index the Ncol field without lying
	m.Ncol++
	
	//loop through the end of each row and insert the data
	var sliceind int = 0 //tracks data index
	for i := m.Ncol-1; i < len(m.V); i += m.Ncol {
		copy(m.V[i+1:], m.V[i:]) //this opens the blank space in the slice
		m.V[i] = data[sliceind]
		sliceind++
	}
}	

//this method vectorizes a matrix by prepending the ncol
func (m *Matrix) Vectorize() (v []float64) {
	//put the Ncol in the slice
	v = append(v, float64(m.Ncol))
	//now add the matrix vector
	v = append(v, m.V...)
	return v
}

//create a function that reads a vectorized matrix into a Go matrix
func AsMatrix(s GoSEXP) (m Matrix) {
	vec := AsFloats(s)
	m.Ncol = int(vec[0])
	m.Nrow = (len(vec) - 1) / m.Ncol
	m.V = vec[1:]
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
	