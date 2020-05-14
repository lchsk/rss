build:
	+$(MAKE) -C api
	+$(MAKE) -C job_runner
	+$(MAKE) -C manage
	+$(MAKE) -C scheduler

unittest:
	+$(MAKE) unittest -C api
	+$(MAKE) unittest -C job_runner
	+$(MAKE) unittest -C manage
	+$(MAKE) unittest -C scheduler

dbtest:
	+$(MAKE) dbtest -C api
	+$(MAKE) dbtest -C job_runner
	+$(MAKE) dbtest -C manage
	+$(MAKE) dbtest -C scheduler

integrationtest:
	+$(MAKE) -C integration_test

.PHONY: build unittest dbtest integrationtest
