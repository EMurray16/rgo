package rfunc

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
