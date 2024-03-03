package requests

import (
	"Project/localElevator/elevator"
	"Project/localElevator/elevio"
	"Project/config"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

type FloorButtonPair struct {
	Floor  int
	Button elevio.ButtonType
}

func requestsAbove(elev elevator.Elevator) bool {
	for f := elev.Floor + 1; f < config.NUM_FLOORS; f++ { //NUM_Floors etc should maybe be in elevio?
		for btn := 0; btn < config.NUM_BUTTONS; btn++ {
			if elev.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(elev elevator.Elevator) bool {
	for f := 0; f < elev.Floor; f++ {
		for btn := 0; btn < config.NUM_BUTTONS; btn++ {
			if elev.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(elev elevator.Elevator) bool {

	for btn := 0; btn < config.NUM_BUTTONS; btn++ {
		if elev.Requests[elev.Floor][btn] {
			return true
		}
	}
	return false
}

func ShouldStop(elev elevator.Elevator) bool {
	switch elev.Dirn {
	case elevio.MD_Down:
		return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !requestsBelow(elev)
	case elevio.MD_Up:
		return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !requestsAbove(elev)
	case elevio.MD_Stop:
		return true
	default:
		return true
	}
}

// Only In Direction Variant
func ShouldClearImmediately(elev elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) bool {

	return elev.Floor == btn_floor &&
		((elev.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
			(elev.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
			elev.Dirn == elevio.MD_Stop ||
			btn_type == elevio.BT_Cab)

}

func ClearAtCurrentFloor(elev elevator.Elevator, ch_LocalElevatorServicedOrder chan []FloorButtonPair) elevator.Elevator {
	localElevatorServicedOrder := make([]FloorButtonPair, 0)
	localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_Cab})
	//elev.Requests[elev.Floor][elevio.BT_Cab] = false
	switch elev.Dirn {
	case elevio.MD_Up:
		if !requestsAbove(elev) && !elev.Requests[elev.Floor][elevio.BT_HallUp] {
			//elev.Requests[elev.Floor][elevio.BT_HallDown] = false
			localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallDown})
		}
		//elev.Requests[elev.Floor][elevio.BT_HallUp] = false
		localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallUp})
		
	case elevio.MD_Down:
		if !requestsBelow(elev) && !elev.Requests[elev.Floor][elevio.BT_HallDown] {
			//elev.Requests[elev.Floor][elevio.BT_HallUp] = false
			localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallUp})
		}
		//elev.Requests[elev.Floor][elevio.BT_HallDown] = false
		localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallDown})
		
	case elevio.MD_Stop:
	default:
		// elev.Requests[elev.Floor][elevio.BT_HallUp] = false
		localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallUp})
		// elev.Requests[elev.Floor][elevio.BT_HallDown] = false
		localElevatorServicedOrder = append(localElevatorServicedOrder, FloorButtonPair{Floor: elev.Floor, Button: elevio.BT_HallDown})
		
	}
	ch_LocalElevatorServicedOrder <- localElevatorServicedOrder
	return elev
}

func RequestsChooseDirection(elev elevator.Elevator) DirnBehaviourPair {
	switch elev.Dirn {
	case elevio.MD_Down:
		if requestsBelow(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else if requestsHere(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_DoorOpen}
		} else if requestsAbove(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	case elevio.MD_Up:
		if requestsAbove(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else if requestsHere(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_DoorOpen}
		} else if requestsBelow(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	case elevio.MD_Stop:
		if requestsHere(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_DoorOpen}
		} else if requestsBelow(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else if requestsAbove(elev) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	default:
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
	}
}
