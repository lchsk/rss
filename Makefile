

all: demo api job_runner manage scheduler

demo:
	@echo "Building demo..."
	@(cd demo && go build)

api:
	@echo "Building api..."
	@(cd api && go build)

job_runner:
	@echo "Building job_runner..."
	@(cd job_runner && go build)

manage:
	@echo "Building manage..."
	@(cd manage && go build)

scheduler:
	@echo "Building scheduler..."
	@(cd scheduler && go build)

.PHONY: demo api job_runner manage scheduler
