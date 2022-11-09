

.PHONY: test-action
test-action: ACT_TYPE ?= edited-checked
test-action:
	act -j pr-comment-master pull_request \
		--env GITHUB_STEP_SUMMARY=/tmp/step-summary \
		-e test/events/pull-request.$(ACT_TYPE).json \
		--container-architecture linux/amd64


.PHONY: test
test:
	go test ./...