package elevator

import (
	"Project/config"
	"Project/localElevator/elevio"
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen                   = 1
	EB_Moving                     = 2
)

func (Behaviour ElevatorBehaviour) String() string {
    names := [...]string{"idle", "doorOpen", "moving"}
    if Behaviour < EB_Idle || Behaviour > EB_Moving {
        return "Unknown"
    }
    return names[Behaviour]
}


type Elevator struct {
	Floor              int
	Dirn               elevio.MotorDirection
	Requests           [config.NUM_FLOORS][config.NUM_BUTTONS]bool
	Behaviour           ElevatorBehaviour
	DoorOpenDuration_s float32
}

func Init() Elevator {
	 var requests [config.NUM_FLOORS][config.NUM_BUTTONS]bool
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			requests[floor][button] = false
		}
	}
	return Elevator{
		Floor:              0,
		Dirn:               elevio.MD_Stop,
		Requests:           requests,
		Behaviour:          EB_Idle,
		DoorOpenDuration_s: config.DoorOpenDuration_s,
	}
}

func PrintElevatorStates(e Elevator) {
	fmt.Println("Elevator state:  ", e.Behaviour, "	Elevator floor: ", e.Floor, "	Elevator direction:  ", e.Dirn)
	for floor := 0; floor < config.NUM_FLOORS; floor++ {

		for button := 0; button < config.NUM_BUTTONS; button++ {
			fmt.Print("	", e.Requests[floor][button])
		}
		fmt.Println("")
	}
	
}
