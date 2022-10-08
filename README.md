# Go-Multi-Microservices-Simulation

<img src="https://user-images.githubusercontent.com/113289516/194726145-f89632ff-ba39-43b8-8b18-508a68669378.jpg" width="800" height="400" />

## Repositoire Breakdown Structure
This project consists of some of the main technical techniques that provide a modern simulation of very high performance for a backend server consisting of three microservices.

### Prerequisites
- Go 
- Docker 
- Docker Compose (Usually, it's installed within Docker)
- Redis (In this project Docker Compose will do the magic)
- RabbitMQ (Docker Compose: Hello, World!)
- Postman
- RedisInsight (Optional)

![compose_swarm](https://user-images.githubusercontent.com/113289516/194728309-865b93ab-fc26-4b41-abce-453b2c1ccf54.png)

With a very simple command line the project will be ready to work:
> docker-compose up -d

### Microservices Description

![index](https://user-images.githubusercontent.com/113289516/194726897-de00d0a9-ef05-487c-8600-258dd16bd9a7.jpg)

Microservices, also known as the microservice architecture, are an architectural style which structures an application as a loosely coupled collection of smaller applications. The microservice architecture allows for the rapid and reliable delivery of large, complex applications. Some of the most common features for a microservice are:
- it is maintainable and testable;
- it is loosely coupled with other parts of the application;
- it  can deployed by itself;
- it is organized around business capabilities;
- it is often owned by a small team.

#### Microservice 1: RESTful API - POST Request (Gin Web Framework)

![background](https://user-images.githubusercontent.com/113289516/194727315-ea6d34f2-1132-4733-895c-b712f30131ee.jpg)

First microservice will start working on port "9090" inside our localhost:
```
err := router.Run("localhost:9090")
```
that microservice works at:
> POST Endpoint: /tweet

its body consist of:
> POST body: { creator: String, body: String}

Then it pushes received information to a RabbitMQ Queue

![vmumd57vnbnq7h1z25xw](https://user-images.githubusercontent.com/113289516/194727616-14c304e0-4ac0-4f6b-8fbb-66ac1b6dc0bb.png)

That microservice return "201 Created" Status if everything going well"
```
context.IndentedJSON(http.StatusCreated, tweet)
```
Otherwise, "400 Bad Request": 
```
context.IndentedJSON(http.StatusBadRequest, gin.H{"response": "bad request"})
```

#### Microservice 2: Message Processor (Preserver)
That microservice subscribes to the queue from RabbitMQ and processes the message. Processing messages means saving the message to Redis.

![workers](https://user-images.githubusercontent.com/113289516/194727952-6faa93eb-f252-4c76-bc5a-9c8a96baedea.jpeg)

#### Microservice 3: RESTful API - GET Request (Gin Web Framework)
Third microservice will start working on port "9091" inside our localhost:
```
err := router.Run("localhost:9091")
```
that microservice works at:
> GET Endpoint: /message/list

that retrieving an array of objects with the creator, and tweet body, content that was stored in the Redis database in chronologically descending order

![1 -VfZ76XK11_1fFUnKLiXSA](https://user-images.githubusercontent.com/113289516/194728218-6a1272a4-34e9-43bb-b517-1e7f6b7c5fd8.png)

Also, that microservice works by passing a parameter to the URL to retrieve the object of a specific creator in chronologically descending order: creator: String
```
router.GET("/tweet/list/:creator", getDataFromRedisDatabaseByParameter)
```
