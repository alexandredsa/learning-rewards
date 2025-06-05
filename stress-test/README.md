# Learning Rewards - Stress Testing Suite

This directory contains a stress testing setup using [Vegeta](https://github.com/tsenart/vegeta), a HTTP load testing tool and library. It's designed to test the performance and reliability of the Learning Rewards API endpoints under load.


## Prerequisites

- Go (for installing Vegeta)
- Make

## Installation

1. Install Vegeta (HTTP load testing tool):
```bash
make install-vegeta
```

## Usage

### Running Stress Tests

The default configuration runs a stress test with:
- Rate: 50 requests per second
- Duration: 10 seconds
- Target: POST requests to `/events` endpoint

To run the default stress test:
```bash
make stress-test
```

### Customizing Test Parameters

You can customize the test parameters using environment variables:

```bash
# Run with custom rate and duration
STRESS_RATE=100 STRESS_DURATION=30s make stress-test
```

Available parameters:
- `STRESS_RATE`: Requests per second (default: 50)
- `STRESS_DURATION`: Test duration (default: 10s)
- `TARGET_FILE`: Path to target definitions (default: targets/events.txt)
- `REPORT_FILE`: Path for binary report (default: results/report.bin)
- `TEXT_REPORT`: Path for text report (default: results/report.txt)

### Test Results

After running a test, you can find the results in:
- Binary report: `results/report.bin`
- Human-readable report: `results/report.txt`

### Cleaning Up

To remove test results:
```bash
make clean-stress-test
```

## Adding New Test Scenarios

1. Create a new JSON payload in `bodies/` for your test case
2. Add a new target definition in `targets/` or modify existing ones
3. Run the tests using the Makefile commands

## Notes

- Ensure the target API server is running before executing tests
- The default target is set to `http://localhost:8081`
- Adjust the target URL in `targets/events.txt` if testing against a different endpoint
