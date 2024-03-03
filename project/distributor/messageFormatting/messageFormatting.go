package messageFormatting

import (
	"Project/config"
	"encoding/json"
	"fmt"
	"strconv"
)

type LocalState struct {
	Behaviour string `json:"behaviour"`
	Direction string `json:"direction"`
	Floor     int    `json:"floor"`
}

func NewLocalState() LocalState {
	return LocalState{
		Behaviour: "idle",
		Direction: "stop",
		Floor:     0,
	}
}

type jsonMessageStruct struct {
	ID             int        `json:"ID"`
	ButtonRequests []int      `json:"globalButtonRequests"`
	LocalState     LocalState `json:"localStates"`
}

type costFuncInputStruct struct {
	HallRequests [][]bool             `json:"hallRequests"`
	States       map[string]ElevState `json:"states"`
}

type ElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

func FormatToJSONMessage(id int, cyclicCounterMatrix [][]interface{}, localState LocalState) []byte {
	buttonRequests := make([]int, 0)
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for buttons := 0; buttons < config.NUM_BUTTONS+config.NUM_ELEVATORS-1; buttons += 2 {
			request := cyclicCounterMatrix[floor][buttons]
			switch request.(type) {
			case int:
				buttonRequests = append(buttonRequests, request.(int))
			default:
				fmt.Println("Error: no type match in cyclicCounterMatrix")
			}
		}
	}
	data := jsonMessageStruct{ID: id, ButtonRequests: buttonRequests, LocalState: localState}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte("")
	}
	// fmt.Println(string(jsonData))
	return jsonData
}

func FormatFromJSONMessage(jsonData []byte) (int, [][]int, LocalState) {
	data := jsonMessageStruct{} //Must be pointer!
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error:", err)
	}
	const buttonsAtFloor = config.NUM_BUTTONS + config.NUM_ELEVATORS - 1
	simplifiedButtonRequests := make([][]int, 0)
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		requestsAtFloor := data.ButtonRequests[floor*buttonsAtFloor : buttonsAtFloor+floor*buttonsAtFloor]
		simplifiedButtonRequests = append(simplifiedButtonRequests, requestsAtFloor[:])
	}
	return data.ID, simplifiedButtonRequests, data.LocalState
}

func FormatToCostFuncInput(cyclicCounterMatrix [][]interface{}, localStates []LocalState, aliveElevators [config.NUM_ELEVATORS]int) []byte {
	hallRequests := make([][]bool, 0)
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		var hallRequestsAtFloor [2]bool
		hallIndex := 0
		for buttons := 0; buttons < (config.NUM_BUTTONS-1)*2; buttons += 2 {
			if cyclicCounterMatrix[floor][buttons] == 1 {
				hallRequestsAtFloor[hallIndex] = true
			} else {
				hallRequestsAtFloor[hallIndex] = false
			}
			hallIndex++
		}
		hallRequests = append(hallRequests, hallRequestsAtFloor[:])
	}
	states := make(map[string]ElevState)
	for ID := 0; ID < len(aliveElevators); ID++ {
		if aliveElevators[ID] == 1 {
			cabRequests := make([]bool, 0)
			for floor := 0; floor < config.NUM_FLOORS; floor++ {
				if cyclicCounterMatrix[floor][(config.NUM_BUTTONS-1+ID)*2] == 1 {
					cabRequests[floor] = true
				} else {
					cabRequests[floor] = false
				}
			}
			elevatorState := ElevState{Behavior: localStates[ID].Behaviour, Floor: localStates[ID].Floor,
				Direction: localStates[ID].Direction, CabRequests: cabRequests}
			IdString := strconv.Itoa(ID)
			states[IdString] = elevatorState
		}
	}
	data := costFuncInputStruct{HallRequests: hallRequests, States: states}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte("")
	}
	// fmt.Println(string(jsonData))
	return jsonData
}

func FormatFromCostFuncOutput(costFuncOutput []byte) {
	// return input to give to elevator. Do we only need the output for this elevator?
	//Arguments : output map[string][][2]bool , thisID ? 
	//Return : hall requests
	//use config.MyID
	//Cost function returns an output variable as a map output := new(map[string][][2]bool). //pointer
	//Access the hall requests by dereferencing the pointer: 
	// val := (*output)["id_1"]
	// fmt.Println(val)


}

//function which returns only output of elevator with ThisID?

/* {
    "id_1" : [[Boolean, Boolean], ...],
    "id_2" : ...
} */
