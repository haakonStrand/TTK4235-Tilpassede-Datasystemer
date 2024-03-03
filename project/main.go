package main

import (
	"Project/config"
	"Project/distributor"
	"Project/distributor/messageFormatting"
	"Project/localElevator/elevio"
	"Project/localElevator/fsm"
	"Project/localElevator/requests"
	"Project/localElevator/timer"
	"Project/network"
	"fmt"
	"flag"
)



func main() {
	elevio.Init("localhost:15657", config.NUM_FLOORS)

	var id int
	flag.IntVar(&id, "id", 0, "id of this peer")
	flag.Parse()
	config.UpdateThisID(id)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ch_ResetDoor := make(chan bool)
	ch_onDoorTimeout := make(chan bool)
	ch_NetworkSendMessage := make(chan []byte)
	ch_NetworkMessageReceived := make(chan []byte)
	ch_LocalElevatorReseivedOrder := make(chan requests.FloorButtonPair,)
	ch_LocalElevatorServicedOrder := make(chan []requests.FloorButtonPair,)
	ch_UpdateLocalElevator := make(chan distributor.LightsOrderPair)
	ch_ElevatorLocalStates := make(chan messageFormatting.LocalState)

	fmt.Println("Elevio init")
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go timer.TimerDoor(config.DoorOpenDuration_s, ch_ResetDoor, ch_onDoorTimeout)
	go fsm.FSM_rutine(drv_floors, drv_buttons, ch_ResetDoor, ch_onDoorTimeout, ch_LocalElevatorServicedOrder, ch_LocalElevatorReseivedOrder, ch_UpdateLocalElevator, ch_ElevatorLocalStates)
	go network.InitNetwork(ch_NetworkSendMessage, ch_NetworkMessageReceived, config.ThisID)
	go distributor.Distributor_routine(ch_NetworkSendMessage, ch_NetworkMessageReceived, 
		ch_LocalElevatorServicedOrder, ch_LocalElevatorReseivedOrder, ch_UpdateLocalElevator, 
		ch_ElevatorLocalStates)
	

	select {}

}
