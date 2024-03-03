
package distributor

import (
	"Project/config"
	"testing"
)

func TestCyclicCounterMatrixInitialize(t *testing.T) {
	cyclicCounterMatrix := cyclicCounterMatrixInitialize()
	PrintCyclicCounterMatrix(cyclicCounterMatrix)
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS*2; button++ {
			if button % 2 == 0 {
				if cyclicCounterMatrix[floor][button] != 0 {
					t.Errorf("Expected 0, got %d", cyclicCounterMatrix[floor][button])
				}
			} else {
				var want [config.NUM_ELEVATORS]int
				if cyclicCounterMatrix[floor][button] != want {
					t.Errorf("Expected [0,0,0], got %v", cyclicCounterMatrix[floor][button])
				}
			}
		}
	}
}

