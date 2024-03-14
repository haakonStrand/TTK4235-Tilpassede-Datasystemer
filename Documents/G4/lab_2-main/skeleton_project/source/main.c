#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <time.h>
#include <unistd.h> //Header file for sleep(). man 3 sleep for details.
#include <pthread.h>
#include "driver/elevio.h"

/**
 * @file main.c
 * @brief The main function operating the logic behind the elevator
 */

typedef struct
{
    int tickets[5][3]; // Ticket matrix containing orders
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

/**
 * @brief function moves elevator to the for it has recieved an order from
 *
 * @param[in] orders
 * @param[in] floor
 */
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

/**
 * @brief this function listens to orders when while sleeping, so no order is ever skipped
 *
 * @param orders
 */
void listenOrders(Orders *orders)
{
    int seconds_to_sleep = 2;
    clock_t start_time = clock();
    while (((clock() - start_time) / CLOCKS_PER_SEC) < seconds_to_sleep)
    {
        for (int f = 0; f < N_FLOORS; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                int btnPressed = elevio_callButton(f, b);
                if (orders->tickets[f][b] != 1 && btnPressed == 1)
                {
                    elevio_buttonLamp(f, b, btnPressed); // Turn on button lamp
                    orders->tickets[f][b] = 1;           // Inserts order
                }
            }
        }
    }
}

/**
 * @brief function that clears the orders, used when the stop button is pressed to remove light and orders.
 *
 * @param[in] orders
 */
void clearOrders(Orders *orders)
{
    for (int f = 0; f < N_FLOORS; f++)
    {
        for (int b = 0; b < N_BUTTONS; b++)
        {
            if (orders->tickets[f][b] == 1)
            {
                elevio_buttonLamp(f, b, 0); // Turns all lights for buttons
                orders->tickets[f][b] = 0;  // Clears all orders
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
            elevio_buttonLamp(f, b, 0); // Turn off lamp
            orders.tickets[f][b] = 0;
        }
    }

    printf("=== Example Program ===\n");

    elevio_motorDirection(DIRN_UP);

    int emergency_stop = 0;
    int last_floor = 0;
    set_default_floor();                  // Elevator goes to first floor by default
    MotorDirection direction = DIRN_STOP; // Create a variable to keep track of last moving state
    MotorDirection last_direction = DIRN_STOP;
    while (1)
    {
        int floor = elevio_floorSensor();

        if (floor != -1)
        {
            last_floor = floor;
            elevio_floorIndicator(floor); // floor light
        }

        // Register ticket to floor
        for (int f = 0; f < N_FLOORS; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                int btnPressed = elevio_callButton(f, b);
                if (orders.tickets[f][b] != 1 && btnPressed == 1)
                {
                    elevio_buttonLamp(f, b, btnPressed); // Turn on elevator lamp
                    orders.tickets[f][b] = 1;            // Collects a ticket
                }
            }
        }

        // Takes inn order and decides motor direction
        for (int f = 0; f < N_FLOORS; f++)
        {
            for (int b = 0; b < N_BUTTONS; b++)
            {
                if (last_floor < f && orders.tickets[f][b] == 1 && !elevio_obstruction())
                {
                    elevio_doorOpenLamp(0); // Turn off lamp when arriving at floor
                    if (direction != DIRN_DOWN)
                    {
                        direction = DIRN_UP;
                        last_direction = DIRN_UP;
                        elevio_motorDirection(DIRN_UP);
                    }
                }
                if (last_floor > f && orders.tickets[f][b] == 1 && !elevio_obstruction())
                {
                    elevio_doorOpenLamp(0); // Turn off lamp when arriving at floor
                    if (direction != DIRN_UP)
                    {
                        direction = DIRN_DOWN;
                        last_direction = DIRN_DOWN;
                        elevio_motorDirection(DIRN_DOWN);
                    }
                }
                if (floor == f && orders.tickets[f][b] == 1)
                {
                    orders.tickets[f][b] = 0;
                    elevio_buttonLamp(f, b, 0); // Turn off lamp
                    elevio_doorOpenLamp(1);     // Turn on lamp when arriving at floor
                    direction = DIRN_STOP;
                    last_direction = DIRN_STOP;
                    elevio_motorDirection(DIRN_STOP);
                    listenOrders(&orders);
                    elevio_doorOpenLamp(0);
                }
                if (last_floor == f && emergency_stop == 1 && orders.tickets[f][b] == 1)
                {
                    if (last_direction == DIRN_UP)
                    {
                        direction = DIRN_DOWN;
                        last_direction = DIRN_DOWN;
                        elevio_motorDirection(DIRN_DOWN);
                        emergency_stop = 0;
                    }
                    else if (last_direction == DIRN_DOWN)
                    {
                        direction = DIRN_UP;
                        last_direction = DIRN_UP;
                        elevio_motorDirection(DIRN_UP);
                        emergency_stop = 0;
                    }
                }
            }
        }

        // go_to_floor(&orders, floor);

        if (elevio_obstruction())
        {
            while (elevio_obstruction())
            {
                for (int f = 0; f < N_FLOORS; f++)
                {
                    for (int b = 0; b < N_BUTTONS; b++)
                    {
                        int btnPressed = elevio_callButton(f, b);
                        if (orders.tickets[f][b] != 1 && btnPressed == 1)
                        {
                            elevio_buttonLamp(f, b, btnPressed); // Turn on elevator lamp
                            orders.tickets[f][b] = 1;            // Collects a ticket
                        }
                    }
                }
                elevio_stopLamp(1);
                elevio_doorOpenLamp(1);
            }
            listenOrders(&orders);
            elevio_stopLamp(0);
            elevio_doorOpenLamp(0);
        }
        else
        {
            elevio_stopLamp(0);
            elevio_doorOpenLamp(0);
        }

        if (elevio_stopButton())
        {
            while (elevio_stopButton())
            {
                direction = DIRN_STOP;
                emergency_stop = 1;
                elevio_motorDirection(DIRN_STOP);
                clearOrders(&orders);
                elevio_stopLamp(1);
                int current = elevio_floorSensor();
                printf("%d", current);
                if (current == 0 || current == 1 || current == 2 || current == 3)
                {
                    elevio_doorOpenLamp(1);
                }
            }

            elevio_doorOpenLamp(0);
            elevio_stopLamp(0);
        }

        nanosleep(&(struct timespec){0, 20 * 1000 * 1000}, NULL);
    }

    return 0;
}
