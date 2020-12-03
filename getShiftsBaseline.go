package AIEngineerChallenge

//a simple, non-optimal algorithm. used for comparison
func getDriverShiftsBaseline(deliveries []*Delivery) []DriverShift {
	if len(deliveries) == 0 {
		return nil
	}
	deliveriesAssigned := make([]bool, len(deliveries))

	shifts := make([]DriverShift, 1) //start with one shift

	//current shift and state
	currentShift := &shifts[0]
	yardLocation := [2]float64{0, 0}
	currentLocation := yardLocation
	currentShiftMinutes := 0.0

	for true {
		//find closest unassigned delivery
		var closestDelivery *Delivery
		var closestDeliveryIdx int
		var minMinutes float64
		for i, d := range deliveries {
			if deliveriesAssigned[i] {
				continue
			}
			minutes := getTravelMinutesBetweenPoints(currentLocation, d.Location)
			if closestDelivery == nil || minutes < minMinutes {
				closestDelivery = d
				closestDeliveryIdx = i
				minMinutes = minutes
			}
		}
		if closestDelivery == nil { //none left; we're done
			return shifts
		} else {
			//can driver get home in time if they do the delivery? if so, assign
			minutesToDoJob := minMinutes + deliveryMinutes
			earliestHomeMinutes := currentShiftMinutes + minutesToDoJob + getTravelMinutesBetweenPoints(closestDelivery.Location, yardLocation)
			if earliestHomeMinutes <= driverShiftLengthMinutes {
				currentShift.Deliveries = append(currentShift.Deliveries, closestDelivery)
				currentShiftMinutes += minutesToDoJob
				currentLocation = closestDelivery.Location
				deliveriesAssigned[closestDeliveryIdx] = true
			} else { //else, don't assign, and get a new driver
				shifts = append(shifts, DriverShift{})
				currentShift = &shifts[len(shifts)-1]
				currentLocation = yardLocation
				currentShiftMinutes = 0.0
			}
		}
	}
	//should never get here
	return shifts
}

//Here's a dumber one
func getDriverShiftsNaive(deliveries []*Delivery) []DriverShift {
	shifts := make([]DriverShift, len(deliveries))
	for i, d := range deliveries {
		shifts[i] = DriverShift{
			Deliveries: []*Delivery{d},
		}
	}
	return shifts
}
