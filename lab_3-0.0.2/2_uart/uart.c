#include "gpio.h"
#include "uart.h"

void uart_init()
{
    // Setting pin 6 (transmit) as output
    GPIO->OUT = (1 << 6);
    // Setting pin 8 (Recieve) as input
    GPIO->IN = (1 << 8);

    // Setting UART pin select to corresponding pins
    UART->PSELTXD = 6;
    UART->PSELRXD = 8;

    // Setting baudrate to 9600
    UART->BAUDRATE = 0x00275000;

    // Disconnect CTS and RTS connections
    UART->PSELCTS = 0xFFFFFFFF;
    UART->PSELRTS = 0xFFFFFFFF;

    // Enable UART
    UART->ENABLE = (4 << 0);

    // Recieve messages enabled
    UART->TASKS_STARTRX = 1;
}

void uart_send(char letter)
{
    UART->TASKS_STARTTX = 1;
    UART->TXD = letter;
    while(!UART->EVENTS_TXDRDY){

    }
    UART->TASKS_STOPTX = 1;
    UART->TASKS_STARTTX = 0;
    UART->EVENTS_TXDRDY = 0;
}

void uart_read()
{
    char temp;
    while(!UART->EVENTS_RXDRDY){

    }
    UART->EVENTS_RXDRDY = 0;
    temp = UART->RXD;
    UART->TASKS_STOPTX = 1;
    return temp;
}
