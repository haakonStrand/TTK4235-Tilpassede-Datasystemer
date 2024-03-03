package fsm

import (
	"Project/localElevator/elevator"
	"Project/localElevator/elevio"
	"Project/localElevator/requests"
	"fmt"
	"Project/distributor"
	"Project/distributor/messageFormatting"
)
func SetElevatorLights(lightsOrderPair distributor.LightsOrderPair, elevatorInstance elevator.Elevator) { //Bedre kodekvalitet å bare ta in lights matrisen?
	elevio.SetFloorIndicator(elevatorInstance.Floor)
	for f := range elevatorInstance.Requests {
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), f, lightsOrderPair.Lights[f][elevio.BT_Cab])
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_HallDown), f, lightsOrderPair.Lights[f][elevio.BT_HallDown])
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_HallUp), f, lightsOrderPair.Lights[f][elevio.BT_HallUp])
	}
}

func elevatorToLocalState(elevatorInstance elevator.Elevator) messageFormatting.LocalState {
	var localState messageFormatting.LocalState
	localState.Behaviour = elevatorInstance.Behaviour.String()
	localState.Floor = elevatorInstance.Floor
	localState.Direction = elevatorInstance.Dirn.String()
	return localState
}

func FSM_rutine(ch_onFloorArrival chan int,
	ch_onRequestButtonPress chan elevio.ButtonEvent,
	ch_ResetDoor chan bool,
	ch_onDoorTimeout chan bool,
	ch_LocalElevatorServicedOrder chan []requests.FloorButtonPair,
	ch_LocalElevatorReceivedOrder chan requests.FloorButtonPair,
	ch_UpdateLocalElevator chan distributor.LightsOrderPair,
	ch_ElevatorLocalStates chan messageFormatting.LocalState){

	elev := elevator.Init()
	elevatorInstance := &elev
	lightsOrderPair := distributor.LightsOrderPair{
		LocalOrders: elevatorInstance.Requests, 
		Lights: elevatorInstance.Requests}  

	elevio.SetMotorDirection(elevio.MD_Down)
	floor := <-ch_onFloorArrival
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevatorInstance.Floor = floor
		

	elevator.PrintElevatorStates(*elevatorInstance)
	
	for {
		select {

		case floor := <-ch_onFloorArrival:
			elevatorInstance.Floor = floor
			switch elevatorInstance.Behaviour {
			case elevator.EB_Moving:
				if requests.ShouldStop(*elevatorInstance) { 
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					requests.ClearAtCurrentFloor(*elevatorInstance, ch_LocalElevatorServicedOrder)
					ch_ResetDoor <- true
					elevatorInstance.Behaviour = elevator.EB_DoorOpen
				}
				SetElevatorLights(lightsOrderPair, *elevatorInstance)
				elevator.PrintElevatorStates(*elevatorInstance)
			default:
				break
			}
			

		case <-ch_onDoorTimeout:
			fmt.Println("Door timer timed out!")
			switch elevatorInstance.Behaviour {
			
			case elevator.EB_DoorOpen:
				var pair requests.DirnBehaviourPair = requests.RequestsChooseDirection(*elevatorInstance)
				elevatorInstance.Behaviour = pair.Behaviour
				elevatorInstance.Dirn = pair.Dirn

				elevio.SetDoorOpenLamp(false)
				elevio.SetMotorDirection(elevatorInstance.Dirn)

				if elevatorInstance.Dirn == elevio.MD_Stop {
					elevatorInstance.Behaviour = elevator.EB_Idle
		
				} else {
					elevatorInstance.Behaviour = elevator.EB_Moving
				}
				elevator.PrintElevatorStates(*elevatorInstance)
			}

		case buttonPress := <-ch_onRequestButtonPress:
			ch_LocalElevatorReceivedOrder <- requests.FloorButtonPair{Floor: buttonPress.Floor, Button: buttonPress.Button}
			
			//Open door emediately if the elevator is already at the floor and the door is open
			if elevatorInstance.Behaviour == elevator.EB_DoorOpen && requests.ShouldClearImmediately(*elevatorInstance, buttonPress.Floor, buttonPress.Button) {
				ch_ResetDoor <- true
			}

		case lightsOrderPair = <-ch_UpdateLocalElevator:
			elevatorInstance.Requests = lightsOrderPair.LocalOrders
			SetElevatorLights(lightsOrderPair, *elevatorInstance)
			elevator.PrintElevatorStates(*elevatorInstance)

			if elevatorInstance.Behaviour == elevator.EB_Idle {
				var pair requests.DirnBehaviourPair = requests.RequestsChooseDirection(*elevatorInstance)
				elevatorInstance.Dirn = pair.Dirn
				elevatorInstance.Behaviour = pair.Behaviour
				
				switch elevatorInstance.Behaviour {

				case elevator.EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					ch_ResetDoor <- true
					requests.ClearAtCurrentFloor(*elevatorInstance, ch_LocalElevatorServicedOrder)

				case elevator.EB_Idle:
					break

				case elevator.EB_Moving:
					elevio.SetMotorDirection(elevatorInstance.Dirn)
					
				}
				
			}

			SetElevatorLights(lightsOrderPair, *elevatorInstance)
			ch_ElevatorLocalStates <- elevatorToLocalState(*elevatorInstance)

			elevator.PrintElevatorStates(*elevatorInstance)
			
		}
	}
}
