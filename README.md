# EnviroStats 
#### *For the course Net Computing*

## Introduction
EnviroStats provides a public API to visualize data throughout the city. Sensors can send their measurements to one of our servers and then the data is visualized on an interactive Google Map. At the moment we generate our own dummy data, for one sensor type, CO² levels. Adding more sensor types as the systems grows larger is easy. The system could be expanded for temperature data, moisture data, noise data, whatever the technology of sensors offers.

Our web application displays the measurement data from our servers on a Google Map like this:

[!DummyPollutionMap](https://github.com/sisisytze/netComputing/blob/master/@images/DummyPollutionMap.png?raw=true)

## Technical README
The system exists out of several components which are shown in the image below. 

### Logical Component Diagram

[!LogicalComponentDiagram](https://github.com/sisisytze/netComputing/blob/master/@images/LogicalComponentDiagram.png?raw=true)

### System Example
[!FourServerExampleSetup](https://github.com/sisisytze/netComputing/blob/master/@images/FourServerExampleSetup.png?raw=true)

In this example we have 1 routing server, 4 data servers (2 pairs), and 15 sensors.

#### Sensors
An amount of n sensors (Python) can send messages using RabbitMQ through sockets. 

When a sensor is introduced to the system it connects to the central routing server, asking which of our servers it should send measurements to. It stores this IP internally and rechecks this IP on a set interval (e.g. daily or weekly). 

*Since we do not have any physical sensors that can send data we have set up a dummyDataGenerator that pretends to a few dozen sensors sending measurements.*

#### Data Servers
Our data servers are the main component of our system, as many different parts run on these servers. On one server we have the following parts:
+ **Java Server** for handling incoming measurements form the message queue (RabbitMQ) and paired syncing of data.
+ **MySQL Database** for storing measurements, and sensor data.
+ **GoLang REST API** for (public) access to the data in the databases accross multiple data servers.
+ **Web Application** to display measurement data using data from the REST API.

#### Routing Server
One centralized routing server that knows all the data servers and also keeps track of pairing. Fulfilling a request for a webpage can be done by any data server, so the routing server just forwards the request to any server. This allows for some load balancing.

Each data server needs to know what data server it is paired with, so that they can exchange their measurement data in an hourly sync. So when a data server is introduced to the system it connects to the central routing server to find out about its pair, which is then stored internally.

### MySQL Database Structure
The table structure of our MySQL database is described in the following image:

[!DatabaseDiagram](https://github.com/sisisytze/netComputing/blob/master/@images/DatabaseDiagram.png?raw=true)
