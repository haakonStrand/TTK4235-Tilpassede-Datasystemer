#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <time.h>
#include "driver/elevio.h"

/**
 * @file main.c
 * @brief The main function operating the logic behind the elevator
 */

struct orders
{
    int tickets[0]; // Ticket vector containing orders
};

/**
 * @brief makes elevator go to default floor (first floor)
 *
 */
void set_default_floor(void)
{
    while (1)
    {
        int floor = elevio_floorSensor();

        if (floor > 0)
        {
            elevio_motorDirection(DIRN_DOWN);
        }
        if (floor == 0)
        {
            elevio_motorDirection(DIRN_STOP);
            return;
        }
    }
}

int main()
{
    elevio_init();

    printf("=== Example Program ===\n");
    printf("Press the stop button on the elevator panel to exit\n");

    elevio_motorDirection(DIRN_UP);

    set_default_floor();

    int floorind = -1;

    while (0)
    {
        int floor = elevio_floorSensor();
        if (floor != floorind && floor != -1)
        {
            printf("%d\n", floor);
        }
        floorind = floor;

        if (floor == 0)
        {
            elevio_motorDirection(DIRN_UP);
        }

        if (floor == N_FLOORS - 1)
        {
            elevio_motorDirection(DIRN_DOWN);
        }

        for (int f = 0; f < N_FLOORS - 1; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                int btnPressed = elevio_callButton(f, b);
                elevio_buttonLamp(f, b, btnPressed);
            }
        }

        if (elevio_obstruction())
        {
            elevio_stopLamp(1);
        }
        else
        {
            elevio_stopLamp(0);
        }

        if (elevio_stopButton())
        {
            elevio_motorDirection(DIRN_STOP);
            break;
        }

        nanosleep(&(struct timespec){0, 20 * 1000 * 1000}, NULL);
    }

    return 0;
}
