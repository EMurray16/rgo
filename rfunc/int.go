package rfunc

func IntMin (InSlice []int) (int) {
	MinValue := InSlice[0]
	for _,val := range InSlice {
		if val < MinValue {
			MinValue = val
		}
	}
	return MinValue
}

func IntMax (InSlice []int) (int) {
	MaxValue := InSlice[0]
	for _,val := range InSlice {
		if val > MaxValue {
			MaxValue = val
		}
	}
	return MaxValue
}

func IntMean (Vec []int) (int) {
	Sum := 0
	Nvals := 0
	for _,val := range Vec {
			Sum += val
			Nvals++
	}
	Mean := Sum / Nvals
	return Mean
}

func IntSum (InSlice []int) (int) {
	Sum := 0
	for _,val := range InSlice {
		Sum += val
	}
	return Sum
}

func IntIn(refint int, checkslice []int) (bool) {
	for _,i := range checkslice {
		if i == refint {
			return true
		}
	}
	return false
}