build:
	@echo "Building job_runner..."
	@go build

unittest:
	@echo "Job runner unittests"

dbtest:
	@echo "Job runner dbtests"

run_integration: build
	INTEGRATION_TEST=1 ./job_runner

.PHONY: build unittest dbtest run_integration
