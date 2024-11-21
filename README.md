# portk
`portk` is a command-line tool designed to help developers manage processes that are running on specific ports. This can be particularly useful when you have multiple services or processes running simultaneously, and you need to quickly identify and terminate processes listening on certain ports.

## Installation
To install `portk`, you can use Go modules:

```bash
go install github.com/gurel/portk
```

## Usage
To find and kill a process that is listening on a specific port, use the following command:

```
portk <port_number>
```

For example, to kill a process listening on port 8080:

```
portk 8000
```

the command will attempt to gracefully shutdown the process with a deadline of 3 seconds before killing it forcefully. to control the graceful shutdown time, you can pass an additional argument after the port number:

```
portk -w 10s 8000
```

If zero if provided as the waiting period, the command will skip gracefull termination and just kill the process.

## Supported Operating Systems
portk is only tested on MacOS, but should also support Linux. If you attempt to use it on other operating systems, it will return an error. Give me a shout if you use it on Linux so i can update this file.

## License
This project is licensed under the MIT License - see the LICENSE file for details.