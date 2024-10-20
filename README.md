[![Coverage Status](https://coveralls.io/repos/github/tiakavousi/codeplayground/badge.svg?branch=main)](https://coveralls.io/github/tiakavousi/codeplayground?branch=main)

# CodePlayground
**CodePlayground** is a web-based multi-language code execution service that allows users to submit code in various programming languages (Python, JavaScript, Bash, Java, C, C++) and receive output in real-time. The backend is built using Go and leverages the Docker SDK and system compilers for handling different programming languages.

## Features
- Supports multiple programming languages: Python, JavaScript, Bash, Java, C, and C++.
- Compilation and Execution: For compiled languages like C, C++, and Java, the service compiles the submitted code and runs the resulting binary.
- Security through isolation: Future support for Docker containerization to safely isolate code execution.
- REST API: Submit code via a simple HTTP POST request and receive the output.

## Languages Supported
- Python
- JavaScript (Node.js)
- Bash
- Java
- C (compiled using gcc)
- C++ (compiled using g++)

## Build
```
docker-compose up -d --build
```

## Test
```
# Browse
http://127.0.0.1:3000/
# Or
$ brew install websocat
$ websocat ws://localhost:8080/execute
{"language":"python","code":"name = input('Enter your name: ')\nprint(f'Hello, {name}!')"}
```

## Tear Down
```
docker-compose down --rmi all
```


## License
Code released under the [MIT License](https://github.com/tiakavousi/codeplayground/blob/main/LICENSE.md).
