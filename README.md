# Worst Choreography Implementation

This repository is created to show a concept to my friend Sacettin

Check the tests especially [./choreography/integration_flow_test.go] to see flow

## Run

```bash
go run cmd/first_step/main.go & go run cmd/second_step/main.go & go run cmd/third_step/main.go
```

## Test

You can run tests with `go test ./...`

## Missing Parts

- Go Map `ok` checks, will panic in any scenario that is not happy path.
- Thread safe maps?
- Validating if request with same Correlation ID is in the inbox of previous step.
