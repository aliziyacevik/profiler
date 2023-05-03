# Profiler

Profiler is a lightweight Go application that works as a reverse proxy to profile incoming HTTP requests and save the request information to a file in JSON format.

## Features

- Profile incoming HTTP requests.
- Save request information to a file in JSON format.
- Works as a reverse proxy to forward requests to the target server.
- Extracts information such as method, URL, headers, IP, browser, and timestamp.
- Configurable using a `.env` file.

## Getting Started

### Prerequisites

- Go 1.16 or higher.

### Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/profiler.git
```

2. Change to the repository's directory
```sh
cd profiler
```

3.Build the application
```sh
go build
```


## Configuration

Create a `.env` file in the repository's directory with the following content:

```sh
PORT=8080
OUT_FILE="profile.json"
OUT_FORMAT="JSON"
TARGET_SERVER="http://localhost:50133"
```
Adjust the values as needed:

- PORT: The port on which the Profiler server will listen.
- OUT_FILE: The path to the output file where request information will be saved.
- OUT_FORMAT: The format in which request information will be saved. Currently, only "JSON" is supported.
- TARGET_SERVER: The URL of the target server to which requests will be forwarded.


## Running the Profiler

```sh
./profiler config.env
```

## Usage

Send an HTTP request to the Profiler server at http://localhost:PORT/your-path. The server will profile the request and save the information to the specified output file, then forward the request to the target server.
