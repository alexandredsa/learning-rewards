
# === Stress Testing ===
STRESS_RATE ?= 20
STRESS_DURATION ?= 5s
TARGET_FILE ?= targets/events.txt
REPORT_FILE ?= results/report.bin
TEXT_REPORT ?= results/report.txt

setup-stress-test-dirs:
	mkdir -p targets
	mkdir -p results

stress-test: setup-stress-test-dirs
	@echo "⚡ Running Vegeta attack at rate=$(STRESS_RATE)/s for duration=$(STRESS_DURATION)"
	vegeta attack -rate=$(STRESS_RATE) -duration=$(STRESS_DURATION) -targets=$(TARGET_FILE) > $(REPORT_FILE)
	vegeta report $(REPORT_FILE) > $(TEXT_REPORT)
	@echo "📊 Report written to $(TEXT_REPORT)"


install-vegeta:
	go install github.com/tsenart/vegeta@latest

clean-stress-test:
	rm -rf results