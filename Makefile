

.PHONY: test-action
test-action: ACT_TYPE ?= edited-checked
test-action:
	act pull_request --env GITHUB_STEP_SUMMARY=/tmp/step-summary -e test/events/pull-request.$(ACT_TYPE).json


.PHONY: test
test:
	go test ./...