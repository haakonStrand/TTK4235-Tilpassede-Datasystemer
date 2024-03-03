Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> We think that with parallelism, events occur completely in parallell, so the can occur exactly at the same time? With concurrency, it can look like events occur at the same time, but we are actually switching quickly between threads to create an illusion.  

What is the difference between a *race condition* and a *data race*? 
> Race conditition is when the outcome of the program is dependent on how quickly different parts are executed and their order. These parts could be something we have no control over(?) and can lead to unexpected behaviour/results. A race condition can occur when two threads attemt to access the same "resource" at the same time. A data race occurs when two threads access the same object without synchronization. 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> Veery roughly, it schedules. We do not know. 


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> In many cases we want to be able to make things occur simultaneously to solve problems. We want to use multiple threads so that we can execute some tasks while another part of the program is waiting for a resource and cannot perform its own task.  

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Fibers do not require thread locking, as they automatically yield to other fibres when they have to wait to continue. There seems to be some built in infrastructure so that only one runs at a time, and the programmer does not have to implement any thread locking. 

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Harder and easier. and both

What do you think is best - *shared variables* or *message passing*?
> Maybe message passing. 


