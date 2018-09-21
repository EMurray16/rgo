//This contains utilities for embarassing paralellization on a set number of corps
package embarrassync

import (
	"sync"
)


//create a struct with a waitgroup for parallelization
type ParCluster struct {
	ProcsWG, AvailWG sync.WaitGroup
	Nprocs, MaxProcs int
	AvailBool bool //this boolean helps deal with latency
}

//create a function that makes the struct
func MakeParCluster(MaxProcs int) ParCluster {
	//create the individual componenets
	wg := sync.WaitGroup{}
	Nprocs := 0
	AvailBool := true
	//now put them together
	return ParCluster{wg, wg, Nprocs, MaxProcs, AvailBool}
}

//create an add and done method for the new struct
func (Clust *ParCluster) Add() {	
	if Clust.AvailBool {
		//we can add a cluster
		Clust.ProcsWG.Add(1)
		Clust.Nprocs++
		
		//see if we need to send a message not to add any more
		if Clust.Nprocs >= Clust.MaxProcs {
			Clust.AvailBool = false
			Clust.AvailWG.Add(1)
		}
		
	} else { 
		//we need to wait for some space to become available and try again
		Clust.AvailWG.Wait()
		Clust.Add()
	}	
}
func (Clust *ParCluster) Done() {
	Clust.ProcsWG.Done()
	Clust.Nprocs--
	if Clust.Nprocs == (Clust.MaxProcs-1) && !Clust.AvailBool {
		//this frees up a processor allocaiton
		Clust.AvailWG.Done()
		Clust.AvailBool = true
	}
}

//create a function that waits on all clusters to be finished
func (Clust *ParCluster) Wait() {
	//wait to make sure all the individual processors are done
	Clust.ProcsWG.Wait()
}