Both the C and Go codes give results which are not zero. The result changes when running the program several times and can be far above or below zero. 
We don't know how quickly each thread operates, and this result could be due to one being able to run much faster than the other. The threads are not waiting for each other or communicating in any way. If they were executed alternately this would not be a problem. 
The problem comes from the mulitple threads, making it possible to both add and subtract at the same time, but the program will only account for one of the actions. 

4) We chose to use mutex. This was because we only had two threads and the thread itself should be the only one to unlock the mutex. 