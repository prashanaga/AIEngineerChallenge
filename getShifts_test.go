package AIEngineerChallenge

import (
	"fmt"
	"testing"
	"time"
	"os"
	"io/ioutil"
	"encoding/json"
)

type SolutionReport struct {
	AllDeliveriesAssignedOnce bool
	AllShiftsMeetConstraints  bool
	NumberOfShifts            int
	TotalLengthOfShifts       float64
	TimeElapsed               time.Duration
}

type solutionFunction func([]*Delivery) []DriverShift

func TestSimple(t *testing.T) {
	deliveries := []*Delivery{
		{
			ID:       1,
			Location: [2]float64{10, 30},
		},
	}
	report := getSolutionReport(deliveries, getDriverShifts)
	reportBaseline := getSolutionReport(deliveries, getDriverShiftsBaseline)
	evaluateSolutionReport(&report, &reportBaseline, "TestSimple", t)
}

func TestAFew(t *testing.T) {
	deliveries := []*Delivery{
		{
			ID:       1,
			Location: [2]float64{10, 30},
		},
		{
			ID:       2,
			Location: [2]float64{-50, 60},
		},
		{
			ID:       3,
			Location: [2]float64{100, 90},
		},
		{
			ID:       4,
			Location: [2]float64{-15, 200},
		},
		{
			ID:       5,
			Location: [2]float64{-40, -40},
		},
	}
	report := getSolutionReport(deliveries, getDriverShifts)
	reportBaseline := getSolutionReport(deliveries, getDriverShiftsBaseline)
	evaluateSolutionReport(&report, &reportBaseline, "TestAFew", t)
}

//1000 deliveries
func TestBig(t *testing.T) {
	deliveries := getDeliveriesFromJSONFile("testData/testBig.json")

	report := getSolutionReport(deliveries, getDriverShifts)
	reportBaseline := getSolutionReport(deliveries, getDriverShiftsBaseline)
	evaluateSolutionReport(&report, &reportBaseline, "TestBig", t)
}

func getDeliveriesFromJSONFile(filePath string) []*Delivery {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var deliveries []*Delivery
	if err := json.Unmarshal(byteValue, &deliveries); err != nil {
		fmt.Println(err)
		return nil
	}

	jsonFile.Close()

	return deliveries
}

func getSolutionReport(deliveries []*Delivery, solution solutionFunction) SolutionReport {
	report := SolutionReport{
		AllDeliveriesAssignedOnce: true,
		AllShiftsMeetConstraints:  true,
	}

	t := time.Now()
	driverShifts := solution(deliveries)
	report.TimeElapsed = time.Since(t)

	report.NumberOfShifts = len(driverShifts)

	//check delivery assignment
	for _, d := range deliveries {
		assignmentCount := 0
		for _, ds := range driverShifts {
			for _, dsd := range ds.Deliveries {
				if d.ID == dsd.ID {
					assignmentCount++
				}
			}
		}
		if assignmentCount != 1 {
			report.AllDeliveriesAssignedOnce = false
			break
		}
	}

	//check shift lengths
	for _, ds := range driverShifts {
		shiftLengthMinutes := getShiftLengthMinutes(&ds)
		if shiftLengthMinutes > driverShiftLengthMinutes {
			report.AllShiftsMeetConstraints = false
		}
		report.TotalLengthOfShifts += shiftLengthMinutes
	}

	return report
}

func evaluateSolutionReport(report *SolutionReport, baselineReport *SolutionReport, testName string, t *testing.T) {
	var invalid bool
	if !report.AllDeliveriesAssignedOnce {
		t.Errorf("Failed %v: Not all deliveries were assigned exactly once", testName)
		invalid = true
	}
	if !report.AllShiftsMeetConstraints {
		t.Errorf("Failed %v: Not all shifts were completed within %v minutes", testName, driverShiftLengthMinutes)
		invalid = true
	}
	if invalid {
		return
	}
	fmt.Printf("Your solution number of shifts: %d\n", report.NumberOfShifts)
	fmt.Printf("Baseline solution number of shifts: %d\n\n", baselineReport.NumberOfShifts)
	fmt.Printf("Your solution total length: %f\n", report.TotalLengthOfShifts)
	fmt.Printf("Baseline solution total length: %f\n\n", baselineReport.TotalLengthOfShifts)
	fmt.Printf("Your solution run time: %v\n", report.TimeElapsed)
	fmt.Printf("Baseline solution run time: %v\n\n", baselineReport.TimeElapsed)
}

func getShiftLengthMinutes(shift *DriverShift) float64 {
	//minutes on location
	accumulatedMinutes := deliveryMinutes * float64(len(shift.Deliveries))

	//start at yard
	yardLocation := [2]float64{0, 0}
	currentLocation := yardLocation
	//go to each point
	for _, d := range shift.Deliveries {
		accumulatedMinutes += getTravelMinutesBetweenPoints(currentLocation, d.Location)
		currentLocation = d.Location
	}
	//return to yard
	accumulatedMinutes += getTravelMinutesBetweenPoints(currentLocation, yardLocation)

	return accumulatedMinutes
}
