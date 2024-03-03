package config

const NUM_FLOORS	int  = 4
const NUM_BUTTONS 	int = 3
const DoorOpenDuration_s float32 = 3.0
const NUM_ELEVATORS int = 3
const AssumedDeadElevatorTime_s float32 = 3.0
var ThisID int = 0

func UpdateThisID(newID int){
	ThisID = newID
}