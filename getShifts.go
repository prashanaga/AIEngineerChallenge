package AIEngineerChallenge
import "bitbucket.org/sjbog/go-dbscan"
import "strconv"
import (
	"math"
	"fmt"
)

//Your task is to write a function that takes a list of Deliveries, and organizes them into a list of DriverShifts

//A Delivery object will have an ID and a Location, which is just a cartesian coordinate (like lat long)
type Delivery struct {
	ID       int
	Location [2]float64
}

//A DriverShift is just a list of what deliveries a driver will do, in order
type DriverShift struct {
	Deliveries []*Delivery
}

//Each delivery itself takes exactly 15 minutes to complete, once a driver has reached the site
const deliveryMinutes = 15

//A driver shift cannot exceed 12 hours
const driverShiftLengthMinutes = 12 * 60

//The solution must meet the following criteria
//1. Each delivery should be assigned once
//2. No DriverShift can take longer than driverShiftLengthMinutes
//You may assume that no single delivery will take longer than a full driver's shift

//The length of a shift is equal to the total time traveled + the total time spent at deliveries
//The travel time in minutes between any two points is just the cartesian distance between them
//	e.g. travel time between points (0,1) and (2,3.5) is sqrt((2-0)^2 + (3.5-1)^2) minutes

//ALL DRIVERS START AND END THEIR SHIFTS AT THE DRIVER YARD, LOCATED AT (0,0). They must return to (0,0) within the shift duration.
//So a driver with a shift of just one delivery at point (20, 30) will have a total shift length of 2*36.055 + 15
//i.e. traveling from (0,0) to (20,30), delivering for 15 minutes, and returning to (0,0)

//For two solutions A and B that both meet the minimum criteria, they are compared as follows
//1. A is better than B if it uses fewer total shifts
//2. If A and B use the same number of shifts, then A is better if the total distance traveled is less

//In general, the function should scale well and return results within a few seconds, even for very large problems (1000s of deliveries)
//For example, if solution B is slightly better than solution A, but B takes 30 seconds to A's 1 second, A is probably preferable

//Implement
func getDriverShifts(deliveries []*Delivery) []DriverShift {
	//your code here!
	if(len(deliveries)==1){
    shiftsFinal1:=[]DriverShift{}
    shiftsFinal1=append(shiftsFinal1,DriverShift{deliveries})
    return shiftsFinal1
    }
	
	var clusterer =dbscan.NewDBSCANClusterer( 100.0, 1 )
	
    var data = []dbscan.ClusterablePoint{}
	
	for _, element := range deliveries { 
	
			
    data=append(data, &dbscan.NamedPoint{Name:strconv.Itoa(element.ID), Point:[]float64{element.Location[0], element.Location[1]}})
    //fmt.Print(index)
    }
	clusterer.MinPts = 1
    clusterer.SetEps(100.0)
    clusterer.AutoSelectDimension = false
    // Set dimension manually
    clusterer.SortDimensionIndex = 1
	
	var result  [][]dbscan.ClusterablePoint=nil
	
	result= clusterer.Cluster(data)
	fmt.Print("**********************************************",data)
	shiftsFinal:=[]DriverShift{}
    for _, element := range result {
       deliveries1 := []*Delivery{}
       for _, element1 := range element {

       // fmt.Print(element1.String()[1]-48,element1.GetPoint(),"  ",index1,"a")
       deliveries1=append(deliveries1,&Delivery{ID:int(element1.String()[1]-48),Location: [2]float64{element1.GetPoint()[0], element1.GetPoint()[1]}})
      }
      // shifts := make([]DriverShift, 1)
 
 //fmt.Println("deliveries1***",*deliveries1[0])
 shifts:=getDriverShiftsBaseline(deliveries1)
 //deliveries1 := []Delivery{}
 for _, shiftElement := range shifts{
     shiftsFinal=append(shiftsFinal,shiftElement)
 }
 
 
}
return shiftsFinal
		   
		   
	
}

//You can test the validity of your solution by running "go test -v ./..." in this repo

//a useful function
func getTravelMinutesBetweenPoints(p0 [2]float64, p1 [2]float64) float64 {
	xDist := p1[0] - p0[0]
	yDist := p1[1] - p0[1]
	return math.Sqrt(xDist*xDist + yDist*yDist)
}
