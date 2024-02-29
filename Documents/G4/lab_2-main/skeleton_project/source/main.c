#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <time.h>
#include "driver/elevio.h"

/**
 * @file main.c
 * @brief The main function operating the logic behind the elevator
 */

typedef struct
{
    int tickets[5][3]; // Ticket vector containing orders
} Orders;

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

void go_to_floor(Orders orders, int floor)
{
    for (int f = 0; f < N_FLOORS; f++)
    {
        for (int b = 0; b < N_BUTTONS; b++)
        {
            if (floor < f && orders.tickets[f][b] == 1)
            {
                elevio_motorDirection(DIRN_UP);
            }
            if (floor > f && orders.tickets[f][b] == 1)
            {
                elevio_motorDirection(DIRN_DOWN);
            }
            if (floor == f && orders.tickets[f][b] == 1)
            {
                printf("%d\n", orders.tickets[f][b]);
                orders.tickets[f][b] = 0;
                elevio_buttonLamp(f, b, 0); // Turn off lamp
                elevio_motorDirection(DIRN_STOP);
            }
        }
    }
}

int main()
{
    elevio_init();

    Orders orders;

    for (int f = 0; f < N_FLOORS; f++)
    {
        for (int b = 0; b < N_BUTTONS; b++)
        {
            orders.tickets[f][b] = 0;
        }
    }

    printf("=== Example Program ===\n");
    printf("Press the stop button on the elevator panel to exit\n");

    elevio_motorDirection(DIRN_UP);

    set_default_floor();

    int floorind = -1;

    while (1)
    {
        int floor = elevio_floorSensor();

        if (floor != -1)
        {
            elevio_floorIndicator(floor);
        }

        floorind = floor;

        if (floor == 0)
        {
            // elevio_motorDirection(DIRN_UP);
        }

        if (floor == N_FLOORS - 1)
        {
            // elevio_motorDirection(DIRN_DOWN);
        }

        for (int f = 0; f < N_FLOORS; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                int btnPressed = elevio_callButton(f, b);
                if (orders.tickets[f][b] != 1 && btnPressed == 1)
                {
                    elevio_buttonLamp(f, b, btnPressed); // Skrur på lampen
                    orders.tickets[f][b] = 1;            // Legger inn bestillingen
                }
            }
        }

        for (int f = 0; f < N_FLOORS; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                if (floor < f && floor != -1 && orders.tickets[f][b] == 1)
                {
                    elevio_motorDirection(DIRN_UP);
                }
                if (floor > f && orders.tickets[f][b] == 1)
                {
                    elevio_motorDirection(DIRN_DOWN);
                }
                if (floor == f && orders.tickets[f][b] == 1)
                {
                    printf("%d\n", orders.tickets[f][b]);
                    orders.tickets[f][b] = 0;
                    elevio_buttonLamp(f, b, 0); // Turn off lamp
                    elevio_motorDirection(DIRN_STOP);
                }
            }
        }

        // go_to_floor(&orders, floor);

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
