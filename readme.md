# Sample Project with Go

## Description
Async Queue Processing is a simple application to process incomming requests asyncronously


## Tools that was used

 - [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream), for additional utilities with NATS.
 - PostgreSQL in order to serve as database


## Architectural Concepts

 - Segregation of Event types as
	 - *Integration event*: A state change that is communicated outside of its bounded context  
	 - *Command*: A request to perform work  
	 - *Query*: A request for some information  
	 - *Reply*: An informational response to either a command or query
