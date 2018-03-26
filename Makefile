# Go Toolchain

# https://github.com/golang/dep
GODEP := dep

# https://godoc.org/golang.org/x/tools/cmd/godoc
GODOC := godoc -http:8080

# https://golang.org/cmd/go/#hdr-List_packages
GOLIST := `go list ./... | grep -v /vendor/`

# gocyclo uses this modified list format
GOLIST_DIR := `go list -f '{{ .Dir }}' ./... | grep -v /vendor/`

# https://golang.org/cmd/gofmt
GOFMT := go fmt $(GOLIST)

# https://golang.org/cmd/vet/
GOVET := go vet $(GOLIST)

# https://github.com/golang/lint
GOLINT := golint $(GOLIST)

# https://github.com/onsi/ginkgo
GOTEST := ginkgo -r --randomizeAllSpecs -noColor -cover && echo

# https://github.com/sozorogami/gover
GOVER := gover

# https://golang.org/cmd/cover/
COVER := go tool cover -func=gover.coverprofile

# https://github.com/fzipp/gocyclo
GOCYCLO := gocyclo -top 10 $(GOLIST_DIR)

# Minimum allowable coverage percentage
COVERAGE_TARGET := 75

# Cleans up any leftover .coverprofile files
.PHONY: clean
clean:
	@find . -type f -name '*.coverprofile' -exec rm {} +

# Runs godoc
.PHONY: doc
doc:
	@$(GODOC)

# Installs build-time dependencies
.PHONY: deps
deps:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/sozorogami/gover
	go get -u github.com/jgautheron/gocyclo
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint

# Installs run-time dependencies
.PHONY: ensure
ensure:
	@$(GODEP) ensure

# Format the code. Will return non-zero exit code if any formatting occurred.
.PHONY: fmt
fmt: format

.PHONY: format
format:
	$(eval FMT_OUT := $(shell $(GOFMT)))
	@[ "$(FMT_OUT)" == "" ] || (echo "$(FMT_OUT)" && exit 1)

# Vet the code. Will return non-zero exit code if any vet rules fail.
.PHONY: vet
vet:
	@$(GOVET)

# Lint the code. Will return non-zero exit code if any lint rules fail.
.PHONY: lint
lint:
	@[ '$(shell $(GOLINT))' == "" ] || ($(GOLINT) && exit 1)

# Returns complexity information on the code.
.PHONY: complexity
complexity:
	@$(GOCYCLO)

# Dummy target necessary for makefile variable expansion reasons
.PHONY: .RUN_TESTS
.RUN_TESTS:
	@$(GOTEST)
	@$(GOVER)

.PHONY: coverage
coverage: fmt vet lint clean .RUN_TESTS

	@$(eval COVERAGE_ACTUAL = $(shell $(COVER) | tee /dev/tty | tail -n 1 | awk '{print $$3}' | sed 's/%//'))
	@echo

	@echo "$(COVERAGE_ACTUAL)" | awk 'BEGIN { if ('$(COVERAGE_ACTUAL)' > '$(COVERAGE_TARGET)') { 	\
		print "[PASS] Coverage meets the required threshold of $(COVERAGE_TARGET)%!"; 				\
		exit 0;																						\
	} else { 																						\
		print "[FAIL] Coverage was below required threshold of $(COVERAGE_TARGET)%!"; 				\
		exit 1; 																					\
	} }'

.PHONY: test
test: coverage
	@echo
	@$(GOCYCLO)
