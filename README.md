
# LRU Cache - Redis and In-Memory Library

## Overview

This project implements an LRU (Least Recently Used) cache with both in-memory and Redis, including TTL (Time-to-Live) functionality.  

The cache operations available are __set__, __get__, __getall__, __delete__, __deleteall__, accessed by POST, GET and DELETE HTTP methods.  

It allows for concurrent get and set operations on the cache as well.  

It also includes benchmarking, testing, and API endpoints implemented with the Gin framework. Postman is used for testing the API.  

## Installation

### Prerequisites

*   Go
*   Redis
*   Postman

### Steps

 -  Clone the repository:
>      git clone https://github.com/pragadeesh-mcw/Go-Mini-Project.git  
>      cd Go-Mini-Project

 -  Install dependencies:  
>   Go-Redis v9  
>   Gin Framework

- Install the dependencies by running:  
>      go mod tidy
	
 -  Start Redis server:
>      redis-server

## Usage

### Running the Application

1.  Build and run the application:

	`go run main.go`

1.  The application will start on the default port `8080`. 

## API Endpoints

### Base URL

http://localhost:8080/cache

### Endpoints

#### Set Key-Value Pair

*   **URL:** `/cache`
*   **Method:** `POST`
*   **Request Body:** `{ key": "your-key", "value": "your-value", "ttl": 60 }`  
>   TTL in seconds
*   **Response:** `{ "message": "Key-Value pair set successfully" }`

#### Get Value by Key

*   **URL:** `/cache/:key`
*   **Method:** `GET`
*   **Response:** `{ "key": "your-key", "value": "your-value" }`

#### Get All Keys
*   **URL:** `/cache/:key`
*   **Method:** `GET`
*   **Response:** `{ "Key-Value Pairs": { key1: value1, key2:value2 } }`

#### Delete Key

*   **URL:** `/cache/:key`
*   **Method:** `DELETE`
*   **Response:** `Key deleted successfully` 

#### Delete All Keys
*   **URL:** `/cache/`
*   **Method:** `DELETE`
*   **Response:** : `All keys deleted successfully` 

The same operations can be performed individually for redis and in-memory cache at 
`/redis/`  and `/inmemory/` respectively.

## Benchmarking
To benchmark the performance of the LRU cache:
1.  Run the benchmark tests:
    `go test -bench . -benchtime=1s -benchmem`
1.  The results will show the performance of different operations in the cache.

## Testing

### Running Tests
To run the tests:  
`cd test`  
`go test` 

## Configuration

Configuration for TTL and max cache size is done via POST.  
The default configuration is:
*   `REDIS_ADDR`: Address of the Redis server (default: `localhost:6379`).
*   `REDIS_PASSWORD`: Password for the Redis server (default: `""`).
*   `REDIS_DB`: Redis database number (default: `0`).
*   `SIZE`: Default size is `3`.
*   `TTL`: Default TTL is `60` seconds.
