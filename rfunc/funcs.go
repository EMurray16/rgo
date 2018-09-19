//This script contains various functions that are migrated from R and may come in handy from time to time
package rfunc

import(
	"sort"
)

//sum applied to booleans
func BoolSum(Vec []bool) (int64) {
	outsum := 0
	for _,val := range Vec {
		if val == true {
			outsum++
		}
	}
	return int64(outsum)
}

//Mean,Median,Max,Min,Sum applied to floats
func FloatMean(Vec []float64) (float64) {
	Sum := float64(0)
	Nvals := 0
	for _,val := range Vec {
		//if val > float64(50) {
			Sum += val
			Nvals++
		//}
	}
	Mean := Sum / float64(Nvals)
	return Mean
}
func FloatMedian(Vec []float64) (float64) {
	//start by sorting it
	sort.Float64s(Vec)
	//Find the indexes for the median
	l := len(Vec)
	var Median float64
	if l <= 1 {
		return Vec[0]
	} else if l%2 == 0 {
		//fmt.Println(l/1, l/2-1, l/2)
		Median = (Vec[l/2 - 1] + Vec[l/2]) / float64(2)
	} else {
		Median = Vec[l/2]
	}
	return Median
}
func FloatMin (InSlice []float64) (float64) {
	MinValue := InSlice[0]
	for _,val := range InSlice {
		if val < MinValue {
			MinValue = val
		}
	}
	return MinValue
}
func FloatMax (InSlice []float64) (float64) {
	MaxValue := InSlice[0]
	for _,val := range InSlice {
		if val > MaxValue {
			MaxValue = val
		}
	}
	return MaxValue
}
func FloatSum (InSlice []float64) (float64) {
	Sum := float64(0)
	for _,val := range InSlice {
		Sum += val
	}
	return Sum
}

//min,max,mean applied to ints
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

//min,max,mean applied to int64s
func Int64Min (InSlice []int64) (int64) {
	MinValue := InSlice[0]
	for _,val := range InSlice {
		if val < MinValue {
			MinValue = val
		}
	}
	return MinValue
}
func Int64Max (InSlice []int64) (int64) {
	MaxValue := InSlice[0]
	for _,val := range InSlice {
		if val > MaxValue {
			MaxValue = val
		}
	}
	return MaxValue
}
func Int64Mean (Vec []int64) (int64) {
	Sum := int64(0)
	Nvals := int64(0)
	for _,val := range Vec {
			Sum += val
			Nvals++
	}
	Mean := Sum / Nvals
	return Mean
}
func Int64Sum (InSlice []int64) (int64) {
	Sum := int64(0)
	for _,val := range InSlice {
		Sum += val
	}
	return Sum
}

//create a function that finds the unique values of a slice
func StringUnique (InSlice []string) ([]string) {
	//find how often each one was used
	Freqs := make(map[string]int, len(InSlice))
	for _,val := range InSlice {
		if val != "" {
			Freqs[val]++
		}
	}
	//make the slice of strings
	StrSlice := []string{}
	for ind,_ := range Freqs {
		StrSlice = append(StrSlice, ind)
	}
	//return it
	return StrSlice
}
//use that function to crate a function that finds teh intersection of 2 slices
func StringIntersect(Slice1 []string, Slice2 []string) ([]string) {
	//get just the unique values of each slice
	Uniques1 := StringUnique(Slice1)
	Uniques2 := StringUnique(Slice2)
	
	//create a map to keep track of how many times a 
	CountMap := make(map[string]int, len(Uniques1)+len(Uniques2))
	
	//loop through each value of Uniques1, adding 1 to the countmap
	for _,str := range Uniques1 {
		CountMap[str]++
	}
	//create a value to keep track of how many are in both
	BothCount := 0
	//do the same for Unqiues2
	for _,str := range Uniques2 {
		CountMap[str]++
		if CountMap[str] == 2 {
			BothCount++
		}
	}
	
	//create a slice of strings
	OutSlice := make([]string, BothCount)
	//fill the slice in wherever CountMap == 2
	ind := 0
	for str,count := range CountMap {
		if count == 2 {
			OutSlice[ind] = str
			ind++
		}
	}
	return OutSlice
}

//create a function that tells us if Refstring is in the CheckSlice
func StringIn (Refstring string, CheckSlice []string) (bool) {
	InBool:= false
	for _,val := range CheckSlice {
		if val == Refstring {
			InBool = true
			break
		}
	}
	
	return InBool
}
//create a function that vectorizes the StringIn_Single operations
func StringIn_Vec (RefSlice []string, CheckSlice []string) ([]bool) {
	InBools := make([]bool, len(RefSlice))
	for ind,refstr := range RefSlice {
		InBools[ind] = StringIn(refstr, CheckSlice)
	}
	return InBools
}

//create a function that finds the number in common
func Ncommon_Uint16(Slice1, Slice2 []uint16) (common int) {
	//compare all elements of both slices, add when common
	for _,i1 := range Slice1 {
		for _,i2 := range Slice2 {
			if i1 == i2 {
				common++
			}
		}
	}
	return common
}