# Mock Server Creator

## Introduction
In the fast-paced world of software development, the ability to quickly mock server responses is invaluable. This not only accelerates the development process but also enhances testing efficiency by allowing developers and testers to focus on the application logic without waiting for the backend services to be up and running. The Mock Server Creator is designed to fulfill this need, offering a swift and straightforward way to set up mock servers. It supports both native installations and Docker images, catering to a wide range of development environments and preferences.

## Features
- **Quick Setup**: Get your mock server running in minutes.
- **Native and Docker Support**: Whether you prefer a native setup or a Docker container, we've got you covered.
- **Configurable**: Easily define your endpoints and responses through a simple configuration file.
- **Cross-Platform**: Compatible with Windows, macOS, and Linux.

## Getting Started

### Prerequisites
- For native setup: Ensure you have Go installed on your system, air is also recommended for hot-reloading.
- For Docker setup: Docker must be installed and running on your machine.

### Installation

#### Native Installation
1. Clone the repository to your local machine.
2. Navigate to the project directory and run `air` to start the app.

#### Docker Installation
1. Pull the Docker image from the repository: `docker pull xdung24/mock-server:latest`

### Configuration

#### Creating a Config File
To define your mock servers and their responses, you need to create a config file. This file should be in JSON format and contain the necessary details about the endpoints you wish to mock.

Example `.env`:
```sh
DATA_FOLDER=./data
USE_FSNOTIFY=false
USE_POLLING=false
POLLING_TIME=15
WEB_ENGINE=fiber
```

# Creating a `settings.json` File

## Overview
To configure your application, you need to create a `settings.json` file that matches the structure defined in the `Setting` struct within the application code. This document guides you through the process of creating this configuration file.

## Structure of `settings.json`

The `settings.json` file should follow this structure:

- `name`: A string representing the name of the setting.
- `host`: The host address.
- `port`: The port number.
- `swaggerEnabled`: Whether to enable swagger-ui for the host (Either openapi.json/openapi.yml/openapi.yaml is required for this option to work)
- `requests`: An array of request objects.
- `headers`: An array of global headers.

### Request Object

Each object in the `requests` array should have the following properties:

- `name`: A string representing the name of the request.
- `method`: The HTTP method (e.g., `GET`, `POST`).
- `path`: The path of the request.
- `responses`: An array of response objects.

### Response Object

Each object in the `responses` array should have the following properties:

- `name`: A string representing the name of the response.
- `code`: The HTTP status code.
- `query`: A query string.
- `headers`: An array of headers specific to this response.
- `filePath`: The path to a file containing the response body.

### Header Object

Both the global headers and the headers within each response should follow this structure:

- `name`: The name of the header.
- `value`: The value of the header.

## Example `settings.json`

Below is an example of a `settings.json` file:

```json
{
  "name": "MyMockServer",
  "host": "localhost",
  "port": 8080,
  "requests": [
    {
      "name": "TestRequest",
      "method": "GET",
      "path": "/api/test",
      "responses": [
        {
          "name": "SuccessResponse",
          "code": 200,
          "query": "type=success",
          "headers": [
            {
              "name": "Content-Type",
              "value": "application/json"
            }
          ],
          "filePath": "responses/success.json"
        }
      ]
    }
  ],
  "headers": [
    {
      "name": "Global-Header",
      "value": "GlobalValue"
    }
  ]
}
