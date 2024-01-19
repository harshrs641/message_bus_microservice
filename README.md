# Microservices Bus with Golang and Kafka

This project implements a microservices bus in Golang using Kafka as a message broker. The microservices bus facilitates communication between different services in a distributed system. It includes three components: `microservices-bus`, `service1`, and `service2`, each serving a specific purpose in the overall architecture.

## Components:

### 1. microservices-bus

- **Description:** The main microservices bus that handles service registration, request handling, and logging of messages.
- **Exposed API:** 
  - `/register`: Register a new service with a given name and URL.
  - `/`: Handle incoming requests and forward them to the appropriate service.
  - `/logs`: Retrieve logs of messages exchanged between services.
- **Dependencies:**
  - Kafka as the message broker.
  
### 2. service1

- **Description:** One of the microservices that exposes an API on port 8081.
- **Purpose:** Serve as an example microservice integrated with the microservices bus.
- **API Endpoint:** `http://localhost:8081`

### 3. service2

- **Description:** Another microservice that exposes an API on port 8082.
- **Purpose:** Demonstrate the interaction between multiple microservices through the microservices bus.
- **API Endpoint:** `http://localhost:8082`

## Setup and Execution:

1. **Docker Compose Setup:**
   - Ensure Docker and Docker Compose are installed on your system.
   - Use the provided `docker-compose.yaml` file for setting up Zookeeper, Kafka, `microservices-bus`, `service1`, and `service2`.

2. **Build and Run:**
   - Execute the following commands in the project directory:
     ```bash
     docker-compose up
     ```

3. **Usage:**
   - Register services using the `/register` endpoint on `microservices-bus`.
   - Send requests to the microservices bus at the root (`/`) to trigger communication with registered services.
   - View logs using the `/logs` endpoint on `microservices-bus`.

## Code Structure:

- **`main.go`:** Entry point of the microservices bus application.
- **`microservices/`:** Directory containing other microservices.
- **`driver/`:** Directory containing the Kafka driver implementation.
- **`log/`:** Directory containing the message log implementation.
- **`event/`:** Directory containing the event struct
- **`integration_test.go`:** Entry point to test the integrations.


## Extending the System:

- To add more microservices, follow the structure of `service1` and `service2`.
- Implement specific functionalities in each microservice and register them with the microservices bus.

## Docker Compose Configuration:

- The provided `docker-compose.yaml` file sets up Zookeeper, Kafka, `microservices-bus`, `service1`, and `service2`.
- Ensure the correct configuration of Kafka environment variables for communication between components.

Feel free to explore, modify, and extend this microservices bus as needed for your specific use case.

## Example Usage

### Service Registration

To register a service with the microservices bus, you can use the following `curl` command:

```bash
curl --location 'http://localhost:8080/register' \
--header 'Content-Type: application/json' \
--data '{
    "name": "service1",
    "url": "http://localhost:8081"
}'
```

This command sends a JSON payload to the `/register` endpoint, providing the service name and URL.

### Requesting Logs

To request logs from the microservices bus after registering a service, you can use the following `curl` command:

```bash
curl --location 'http://localhost:8080/logs' \
--header 'X-Service-Name: service1'
```

This command specifies the service name using the `X-Service-Name` header and retrieves logs for the specified service from the `/logs` endpoint.

