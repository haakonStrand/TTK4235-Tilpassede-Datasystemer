/**
 * @file elevio.c
 * @brief elevio.c includes functions for operating the elevator, it includes initialization, motordirection
 * buttonlamp, floorindicator, opendoorlamp, stoplamp, osv.
 */

#include <assert.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <stdio.h>
#include <pthread.h>

#include "elevio.h"
#include "con_load.h"

static int sockfd;
static pthread_mutex_t sockmtx;

/**
 * @brief Initializes the elevator by setting the ip, and the port that should be used. It also loads the
 * elevio.con file that is made from the make file.
 *
 */
void elevio_init(void)
{
    char ip[16] = "localhost";
    char port[8] = "15657";
    con_load("elevio.con",
             con_val("com_ip", ip, "%s")
                 con_val("com_port", port, "%s"))

        pthread_mutex_init(&sockmtx, NULL);

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    assert(sockfd != -1 && "Unable to set up socket");

    struct addrinfo hints = {
        .ai_family = AF_INET,
        .ai_socktype = SOCK_STREAM,
        .ai_protocol = IPPROTO_TCP,
    };
    struct addrinfo *res;
    getaddrinfo(ip, port, &hints, &res);

    int fail = connect(sockfd, res->ai_addr, res->ai_addrlen);
    assert(fail == 0 && "Unable to connect to elevator server");

    freeaddrinfo(res);

    send(sockfd, (char[4]){0}, 4, 0);
}

/**
 * @brief The function sets the motor direction
 *
 * @param[in] dirn from the  enum struct @c MotorDirection , the function then sets the torque direction
 *  based upon input
 */
void elevio_motorDirection(MotorDirection dirn)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){1, dirn}, 4, 0);
    pthread_mutex_unlock(&sockmtx);
}

/**
 * @brief turns on button lamp when pressed
 *
 * @param[in] floor
 * @param[in] button
 * @param[in] value
 */
void elevio_buttonLamp(int floor, ButtonType button, int value)
{
    assert(floor >= 0);
    assert(floor < N_FLOORS);
    assert(button >= 0);
    assert(button < N_BUTTONS);

    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){2, button, floor, value}, 4, 0);
    pthread_mutex_unlock(&sockmtx);
}

/**
 * @brief indicates what floor the elevator is on by turning on an LED on the device.
 *
 * @param[in] floor
 */
void elevio_floorIndicator(int floor)
{
    assert(floor >= 0);
    assert(floor < N_FLOORS);

    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){3, floor}, 4, 0);
    pthread_mutex_unlock(&sockmtx);
}

/**
 * @brief Indicates that the door is open
 *
 * @param[in] value
 */
void elevio_doorOpenLamp(int value)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){4, value}, 4, 0);
    pthread_mutex_unlock(&sockmtx);
}

/**
 * @brief lamp on device that indicates that the elevator has come to a halt.
 *
 * @param[in] value
 */
void elevio_stopLamp(int value)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){5, value}, 4, 0);
    pthread_mutex_unlock(&sockmtx);
}

/**
 * @brief collects information from the device, and returns 1 if the specific button is currently
 * being pressed
 *
 * @param[in] floor
 * @param[in] button
 * @return int
 */
int elevio_callButton(int floor, ButtonType button)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){6, button, floor}, 4, 0);
    char buf[4];
    recv(sockfd, buf, 4, 0);
    pthread_mutex_unlock(&sockmtx);
    return buf[1];
}

/**
 * @brief sensor reading the floors, returns integer value corresponding to
 * the floor the elevator is at.
 *
 * @return int
 */
int elevio_floorSensor(void)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){7}, 4, 0);
    char buf[4];
    recv(sockfd, buf, 4, 0);
    pthread_mutex_unlock(&sockmtx);
    return buf[1] ? buf[2] : -1;
}

/**
 * @brief listens wether or not the red stop button is being pressed,
 * return a value of 1 or 0 depending on if is being pressed. Used to stop elevator.
 *
 * @return int
 */
int elevio_stopButton(void)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){8}, 4, 0);
    char buf[4];
    recv(sockfd, buf, 4, 0);
    pthread_mutex_unlock(&sockmtx);
    return buf[1];
}

/**
 * @brief function checks if the obstruction lever is on or off,
 * returns 1 if active, 0 if not active.
 *
 * @return int
 */
int elevio_obstruction(void)
{
    pthread_mutex_lock(&sockmtx);
    send(sockfd, (char[4]){9}, 4, 0);
    char buf[4];
    recv(sockfd, buf, 4, 0);
    pthread_mutex_unlock(&sockmtx);
    return buf[1];
}
