package distributor

import (
	"Project/config"
	"Project/distributor/messageFormatting"
	"Project/localElevator/requests"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)



type LightsOrderPair struct {
	LocalOrders [config.NUM_FLOORS][config.NUM_BUTTONS]bool
	Lights [config.NUM_FLOORS][config.NUM_BUTTONS]bool
}
 

func cyclicCounterMatrixInitialize() ([][]interface{}){
	var cyclicCounterMatrix [][]interface{}
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		var floorArray []interface{}
		for button := 0; button < config.NUM_BUTTONS-1+config.NUM_ELEVATORS; button++ {
			floorArray = append(floorArray, 0)
			floorArray = append(floorArray, [config.NUM_ELEVATORS]int{})
		}
		cyclicCounterMatrix = append(cyclicCounterMatrix, floorArray)
	}
	return cyclicCounterMatrix
}

func PrintCyclicCounterMatrix(cyclicCounterMatrix [][]interface{}) {
	for floor := 0; floor < len(cyclicCounterMatrix); floor++ {
		for button := 0; button < len(cyclicCounterMatrix[floor]); button++ {
			fmt.Print(cyclicCounterMatrix[floor][button], " ")
		}
		fmt.Println()
	}
}

func recalculateLocalOrders(cyclicCounterMatrix [][]interface{}, 
	localStates []messageFormatting.LocalState, 
	aliveElevators [config.NUM_ELEVATORS]int)([config.NUM_FLOORS][config.NUM_BUTTONS]bool) {
	
	costFunctionInput := messageFormatting.FormatToCostFuncInput(cyclicCounterMatrix, localStates, aliveElevators)

	ret, err := exec.Command("./cost_fns/hall_request_assigner/hall_request_assigner", "-i", string(costFunctionInput)).CombinedOutput()
	
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return [4][3]bool{}
	}

	output := new(map[string][][2]bool)
	
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return [4][3]bool{}
	}

	outputValue := *output
	
	IdString := strconv.Itoa(config.ThisID)
	newestLocalOrderMatrix := [config.NUM_FLOORS][config.NUM_BUTTONS]bool{} //A matrix that will be filled ith the new local orders

	for floor := 0; floor < config.NUM_FLOORS; floor++ {
	//Get the hall requests from the cost function
		for button := 0; button < config.NUM_BUTTONS-1; button++ { 
				newestLocalOrderMatrix[floor][button] = outputValue[IdString][floor][button]
			}
				//Get the cab requests from the CyclicCounter matrix
			if cyclicCounterMatrix[floor][(config.NUM_BUTTONS-1+ config.ThisID)*2] == 2 { 
				newestLocalOrderMatrix[floor][config.NUM_BUTTONS-1] = true
			}else{
			newestLocalOrderMatrix[floor][config.NUM_BUTTONS-1] = false
		}
	}

	return newestLocalOrderMatrix
}

func makeLightsMatrix(cyclicCounterMatrix [][]interface{}) [config.NUM_FLOORS][config.NUM_BUTTONS]bool {
	var lightsMatrix [config.NUM_FLOORS][config.NUM_BUTTONS]bool
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < (config.NUM_BUTTONS-1)*2; button+= 2{
			if cyclicCounterMatrix[floor][button] == 2 {
				lightsMatrix[floor][button] = true
			} else {
				lightsMatrix[floor][button] = false
			}
		}
		if cyclicCounterMatrix[floor][(config.NUM_BUTTONS-1+config.ThisID)*2] == 2 {
			lightsMatrix[floor][config.NUM_BUTTONS-1] = true
		} else {
			lightsMatrix[floor][config.NUM_BUTTONS-1] = false
		}
	}
	return lightsMatrix
}

func Distributor_routine(ch_NetworkSendMessage chan []byte, 
	ch_NetworkMessageReceived chan []byte,
	ch_LocalElevatorServicedOrder chan []requests.FloorButtonPair,
	ch_LocalElevatorReseivedOrder chan requests.FloorButtonPair,
	ch_UpdateLocalElevator chan LightsOrderPair,
	ch_ElevatorLocalStates chan messageFormatting.LocalState) {

	
// Initialize aliveTimers that report a dead elevator 

	ch_ElevatorHasDied := make(chan int)
	var aliveTimers []time.Timer
	for i := 0; i < config.NUM_ELEVATORS; i++ {
		timer := time.NewTimer(time.Duration(config.AssumedDeadElevatorTime_s) * time.Second)
		aliveTimers = append(aliveTimers, *timer)
		go func(timerIndex int, ch_ElevatorHasDied chan int) {
			<-aliveTimers[timerIndex].C
			ch_ElevatorHasDied <- timerIndex
		}(i, ch_ElevatorHasDied)
	}
	aliveTimers[config.ThisID].Stop()

// Initialize distributor matrices and variables
	
	cyclicCounterMatrix := cyclicCounterMatrixInitialize()
	newestLocalOrderMatrix := [config.NUM_FLOORS][config.NUM_BUTTONS]bool{}
	newestLightsMatrix := [config.NUM_FLOORS][config.NUM_BUTTONS]bool{}
	
	var aliveElevators [config.NUM_ELEVATORS]int
	aliveElevators[config.ThisID] = 1

	localStates := make([]messageFormatting.LocalState, config.NUM_ELEVATORS)



// Main distributor loop

	for {
		select {
		case message := <- ch_NetworkMessageReceived:
			
			// Extract data from the message
			id, recievedMatrix, recievedLocalState := messageFormatting.FormatFromJSONMessage(message)
			
			if id == config.ThisID {
				break
			}

			// Reset timer for elevator and set the elevator to be allive
			aliveTimers[id].Reset(time.Duration(config.AssumedDeadElevatorTime_s) * time.Second)
			aliveElevators[id] = 1
			
			// Update the cyclic counter matrix
			for floor := 0; floor < config.NUM_FLOORS; floor++ {
				for button := 0; button < (config.NUM_BUTTONS+config.NUM_ELEVATORS-1)*2; button+= 2{
					
					if recievedMatrix[floor][button] == 1 && cyclicCounterMatrix[floor][button] == 1 { 
						element :=cyclicCounterMatrix[floor][button+1] 
						switch heardFrom:=element.(type) {						
						case [config.NUM_ELEVATORS]int:							
							
							heardFrom[id]=1
							
							for i :=0; i < len(aliveElevators); i++ {
								if aliveElevators[i] == 0 {
									heardFrom[i]=0
								}
							}
							
							if aliveElevators == heardFrom{
								cyclicCounterMatrix[floor][button] = 2								
								for i := 0; i < len(heardFrom); i++ {
									heardFrom[i] = 0
								}
							}

							cyclicCounterMatrix[floor][button+1] = heardFrom
						
						default:
							break
						}
					}

					if recievedMatrix[floor][button] == 1 && cyclicCounterMatrix[floor][button] == 0 {							
						cyclicCounterMatrix[floor][button] = recievedMatrix[floor][button]
						cyclicCounterMatrix[floor][button] = recievedMatrix[floor][button]
						element :=cyclicCounterMatrix[floor][button+1] 
						switch heardFrom:=element.(type) {
						case [config.NUM_ELEVATORS-1]int:
							heardFrom[id]=1
							heardFrom[config.ThisID]=1
							cyclicCounterMatrix[floor][button+1] = heardFrom
						
						default:
							break
						}
					}

					if recievedMatrix[floor][button] == 2 && cyclicCounterMatrix[floor][button] == 1 {
						cyclicCounterMatrix[floor][button] = recievedMatrix[floor][button]
						element :=cyclicCounterMatrix[floor][button+1] 
						switch heardFrom:=element.(type) {
						case [config.NUM_ELEVATORS-1]int:
							for i := 0; i < len(heardFrom); i++ {
								heardFrom[i] = 0
							}
							heardFrom[id]=1
							cyclicCounterMatrix[floor][button+1] = heardFrom	
						default:
							break
						}
					}
					
					if recievedMatrix[floor][button] == 0 && cyclicCounterMatrix[floor][button] == 2 {
						cyclicCounterMatrix[floor][button] = recievedMatrix[floor][button]
					}
				}
			}

			// Update the relevant local state

			localStates[id] = recievedLocalState

			//Format message that should be sent to the network

			networkMessageToSend := messageFormatting.FormatToJSONMessage(config.ThisID, cyclicCounterMatrix, localStates[config.ThisID])

			// Send the message to the network
			ch_NetworkSendMessage <- networkMessageToSend

			//Recalculate the local orders
			
			newestLocalOrderMatrix = recalculateLocalOrders(cyclicCounterMatrix, localStates, aliveElevators)
			
			newestLightsMatrix = makeLightsMatrix(cyclicCounterMatrix)

			newLightsOrderPair := LightsOrderPair{newestLocalOrderMatrix, newestLightsMatrix}
			
			// Send the orders to the local elevator
			ch_UpdateLocalElevator <- newLightsOrderPair
		
		
		case id :=<- ch_ElevatorHasDied:
			aliveElevators[id] = 0

			newestLocalOrderMatrix = recalculateLocalOrders(cyclicCounterMatrix, localStates, aliveElevators)
			
			newestLightsMatrix = makeLightsMatrix(cyclicCounterMatrix)

			newLightsOrderPair := LightsOrderPair{newestLocalOrderMatrix, newestLightsMatrix}

			ch_UpdateLocalElevator <- newLightsOrderPair
			


		case order :=<- ch_LocalElevatorServicedOrder:
			for i := 0; i < len(order); i++ {
				cyclicCounterMatrix[order[i].Floor][order[i].Button] = 0
				newestLocalOrderMatrix[order[i].Floor][order[i].Button] = false
				newestLightsMatrix[order[i].Floor][order[i].Button] = false	
			}
			
			newLightsOrderPair := LightsOrderPair{newestLocalOrderMatrix, newestLightsMatrix}
			ch_UpdateLocalElevator <- newLightsOrderPair
      

			networkMessageToSend := messageFormatting.FormatToJSONMessage(config.ThisID, cyclicCounterMatrix, localStates[config.ThisID])
			ch_NetworkSendMessage <- networkMessageToSend

		
		case order :=<- ch_LocalElevatorReseivedOrder:

			cyclicCounterMatrix[order.Floor][order.Button*2] = 1
			
			element :=cyclicCounterMatrix[order.Floor][order.Button*2+1]
			switch v := element.(type) {
			case [config.NUM_ELEVATORS]int:
				v[config.ThisID] = 1
				cyclicCounterMatrix[order.Floor][order.Button*2+1] = v
			
			default:
				break
			}

			networkMessageToSend := messageFormatting.FormatToJSONMessage(config.ThisID, cyclicCounterMatrix, localStates[config.ThisID])
			ch_NetworkSendMessage <- networkMessageToSend
		
		case localState :=<- ch_ElevatorLocalStates:
			localStates[config.ThisID] = localState
		}
	}


}