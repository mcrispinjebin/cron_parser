# Cron Parser

The Cron Parser command-line application validates the given cron string and processes it to display the execution units in each segment.

[Go](https://go.dev/)  

---

**Contents**

1. [Setup](#setup)


---

### Setup ###

1. Install Golang and ensure Go project can be run in the system
1. Change the working directory to cron parser using `cd cron_parser/` command
1. Use the command `go mod vendor` and `go mod tidy` to download the dependencies
1. Use `go run main.go "$cron_string"` to run the application
1. Use `go test ./...` to run the unit test cases

---

