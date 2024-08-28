# Passing Notes

## Overview

This project serves as an example of how to define and assign values to metrics using the Prometheus library within a Go application. It provides a practical demonstration of integrating Prometheus to monitor and collect metrics, making it a valuable reference for developers looking to implement similar functionality in their own Go projects.

## Run Application

### DevContainer (Recommended)

If you choose this method, you will need:

- Docker
- An IDE that supports DevContainers (e.g., VS Code, GoLand)

When the DevContainer starts, it will automatically provision instances of Grafana, Prometheus, and Node-Exporter.

- Port Forwarding: The DevContainer will forward requests from port 8181 on your local machine to port 8181 inside the container. This is the port where your application is configured to listen.
- Prometheus Access: You can access the Prometheus application at localhost:9090.
- Grafana Access: You can access the Grafana dashboard at localhost:3000. Use the following credentials to log in:
    - Username: admin
    - Password: grafana

### Docker-Compose

If you're not interested in the debugging process, you can create all the necessary components by running the command below. This will set up 3 instances of application and 1 ngnix component, as well as Grafana, Prometheus, and 

```
docker compose -f .devcontainer/docker-compose.yml -f .devcontainer/docker-compose.infra.yml up
```

Node-Exporter.
- Prometheus Access: You can access the Prometheus application at localhost:9090.
- Grafana Access: You can access the Grafana dashboard at localhost:3000. Use the following credentials to log in:
    - Username: admin
    - Password: grafana
- Nginx Access: You can access to your application over ngnix that uses localhost:8080 to forward your requests.

## Metric Generation

- **`/index` Endpoint**: You can retrieve the current time via the `/index` endpoint. The primary purpose of this endpoint is to track the change in request counts, which is reflected in the 'Op App - Request Counter' visualization.

- **`/echo/:statusCode` Endpoint**: This endpoint becomes available as soon as the application is running. When you provide a status code in the URL, the endpoint will respond with that exact HTTP status code. This allows you to supply data to the 'Op App - Status Code Counter' metric visualization within the custom dashboard.

- **`Config Tracking`**: When running the application within the DevContainer, you'll have the ability to interact with the data on the 'Op App - Config Tracker' metric visualizer. Any changes made to the `config.yaml` file during the application's runtime will be monitored by the application. A background service periodically updates the Prometheus metric values based on the information in the config file.
