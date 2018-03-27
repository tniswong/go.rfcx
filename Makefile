# Cleans up any leftover .coverprofile files
.PHONY: clean
clean:
	@find . -type f -name '*.coverprofile' -exec rm {} +

# Runs godoc
.PHONY: doc
doc:
	$(eval PORT := 6060)
	@echo Starting godoc server at http://localhost:$(PORT). Enter Ctrl-C to stop.
	@godoc -http :$(PORT)

# Installs build-time dependencies
.PHONY: deps
deps:
	@command -v ginkgo &>/dev/null 	|| go get -u github.com/onsi/ginkgo/ginkgo
	@command -v gover &>/dev/null 	|| go get -u github.com/sozorogami/gover
	@command -v gocyclo &>/dev/null || go get -u github.com/fzipp/gocyclo
	@command -v dep &>/dev/null 	|| go get -u github.com/golang/dep/cmd/dep
	@command -v golint &>/dev/null 	|| go get -u golang.org/x/lint/golint

# Installs run-time dependencies using https://github.com/golang/dep
.PHONY: ensure
ensure:
	@dep ensure

# Alias for format.
.PHONY: fmt
fmt: format

# Format the code
.PHONY: format
format:
	@go fmt `go list ./... | grep -v /vendor/`

# Check the format of the code. Will return non-zero exit code if any formatting is needed. Uses
# https://golang.org/cmd/gofmt/
.PHONY: .CHECK_FORMAT
.CHECK_FORMAT:

	@echo "Formatting Results:"
	@echo "-------------------"

	$(eval FMT_OUT := $(shell gofmt -l `go list -f '{{ .Dir }}' ./... | grep -v /vendor/`))

	@(																			\
		[ "$(FMT_OUT)" == "" ]  												\
		&& echo "[PASS] No unformatted code detected!"							\
		&& echo																	\
	) || ( 																		\
		echo "$(FMT_OUT)" 	 													\
		&& echo "[FAIL] Unformatted code detected! You need to run 'make fmt'!" \
		&& exit 1 																\
	)

# Vet the code. Will return non-zero exit code if any vet rules fail.
.PHONY: vet
vet:
	@echo "Vetting Results:"
	@echo "----------------"

	@( 																				\
		go vet `go list ./... | grep -v /vendor/` 									\
		&& echo "[PASS] No vetting errors detected!"								\
		&& echo																		\
	) || (																			\
		echo																		\
		&& echo "[FAIL] Vetting errors detected! Fix them ALL before you proceed!"	\
		&& exit 1																	\
	)

# Dummy target for console formatting. Fixes console display ordering issues because we're using tee to /dev/tty
.PHONY: .PRE_LINT
.PRE_LINT:
	@echo "Linting Results:"
	@echo "----------------"

# Lint the code. Will return non-zero exit code if any lint rules fail. Uses https://github.com/golang/lint
.PHONY: lint
lint: .PRE_LINT

	$(eval GOLINT_OUT := $(shell golint `go list ./... | grep -v /vendor/` | tee /dev/tty))

	@(																			\
		[ '$(GOLINT_OUT)' == "" ]												\
		&& echo "[PASS] No lint errors detected!"								\
		&& echo																	\
	) || (																		\
		echo "[FAIL] Lint errors detected! Fix them ALL befor eyou proceed!"	\
		&& exit 1																\
	)

# Returns cyclomatic complexity information on the code. Uses https://github.com/fzipp/gocyclo
.PHONY: complexity
complexity:

	@$(eval COMPLEXITY_MAXIMUM := 25)
	@$(eval COMPLEXITY_ACTUAL := $(shell										\
		gocyclo -top 10 `go list -f '{{ .Dir }}' ./... | grep -v /vendor/` |	\
		head -n 1 |																\
		awk '{print $$1}'														\
	))

	@echo "Cyclomatic Complexity Report - Top Ten Most Complex:"
	@echo "----------------------------------------------------"
	@gocyclo -top 10 `go list -f '{{ .Dir }}' ./... | grep -v /vendor/`

	@echo

	@echo "$(COMPLEXITY_ACTUAL)" | awk 'BEGIN { if ('$(COMPLEXITY_ACTUAL)' <= '$(COMPLEXITY_MAXIMUM)') {	\
		print "[PASS] The most complex code was below the maximum threshold of $(COMPLEXITY_MAXIMUM)!";		\
		exit 0;																								\
	} else {																								\
		print "[FAIL] The most complex code was above the maximum threshold of $(COMPLEXITY_MAXIMUM)!";		\
		exit 1;																								\
	} }'

# Dummy target necessary for makefile variable expansion reasons. Ginkgo (https://github.com/onsi/ginkgo) is our test
# runner. Ginkgo leaves *.coverprofile files separated into each package, so we use gover
# (https://github.com/sozorogami/gover) to merge them all together
.PHONY: .RUN_TESTS
.RUN_TESTS:
	@echo "Test Results:"
	@echo "-------------"
	@ginkgo -r --randomizeAllSpecs -noColor -cover && gover && echo

# Dummy target for console formatting. Fixes console display ordering issues because we're using tee to /dev/tty
.PHONY: .PRE_COVERAGE
.PRE_COVERAGE:
	@echo "Coverage Report:"
	@echo "----------------"

# Coverage Report using https://golang.org/cmd/cover/.
.PHONY: coverage
coverage: .CHECK_FORMAT vet lint clean .RUN_TESTS .PRE_COVERAGE

	@$(eval COVERAGE_MINIMUM := 80)
	@$(eval COVERAGE_ACTUAL := $(shell					\
		go tool cover -func=gover.coverprofile |		\
		tee /dev/tty | 									\
		tail -n 1 | 									\
		awk '{print $$3}' | 							\
		sed 's/%//'										\
	))

	@echo

	@echo "$(COVERAGE_ACTUAL)" | awk 'BEGIN { if ('$(COVERAGE_ACTUAL)' > '$(COVERAGE_MINIMUM)') { 	\
		print "[PASS] Coverage meets the required threshold of at least $(COVERAGE_MINIMUM)%!"; 	\
		exit 0;																						\
	} else { 																						\
		print "[FAIL] Coverage was below required threshold of at least $(COVERAGE_MINIMUM)%!"; 	\
		exit 1; 																					\
	} }'

	@echo

# Test Report using https://onsi.github.io/ginkgo/
.PHONY: test
test: coverage complexity
