#!/bin/bash

# Build the Docker image
docker build -t xdung24/mock-servers:1.0 -t xdung24/mock-servers:latest .

# Push the Docker image
docker push xdung24/mock-servers:1.0 
docker push xdung24/mock-servers:latest