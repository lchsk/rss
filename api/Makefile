build:
	@echo "Building api..."
	@go build

unittest:
	@echo "Api unittests"
	@go test -v

dbtest:
	@echo "Api dbtests"
	@go test -v --tags database

run_integration: build
	INTEGRATION_TEST=1 ./api

.PHONY: build unittest dbtest run_integration
