# Live pprof is a Golang library to generate and use live profiles.
# Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
# Live pprof homepage: https://github.com/ufoot/livepprof
# Contact author: ufoot@ufoot.org

.PHONY: all
.PHONY: check
.PHONY: build
.PHONY: test
.PHONY: verbose
.PHONY: race
.PHONY: bench
.PHONY: coverage
.PHONY: install
.PHONY: generate
.PHONY: vet
.PHONY: lint

GITHUB_LIVEPPROF=github.com/ufoot/livepprof
PACKAGES=\
	. \
	objfile \
	collector \
	collector/cpu \
	collector/heap \
	cmd/livepprofdemo

# Default task, run regularly when developping.
all: generate fmt vet build
	@echo "$$(date): all done"

# Typically run by CI (eg, Travis).
check:
	for i in $(PACKAGES) ; do go test -o tmp.test -short $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): check done"

# Build all packages.
build:
	for i in $(PACKAGES) ; do go build $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): build done"

# Test all packages.
test:
	for i in $(PACKAGES) ; do go test -o tmp.test $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): test done"

# Test all packages in verbose mode.
verbose:
	for i in $(PACKAGES) ; do go test -o tmp.test -v $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): verbose done"

# Test all packages for race conditions.
race:
	for i in $(PACKAGES) ; do go test -o tmp.test -short -race $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): race done"

# Benchmark all packages.
bench:
	for i in $(PACKAGES) ; do go test -o tmp.test -short -run=NONE -bench . $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): bench done"

# Get test coverage.
coverage:
	for i in $(PACKAGES) ; do go test -o tmp.test -cover $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): coverage done"

# Format all packages.
fmt:
	for i in $(PACKAGES) ; do go fmt $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): fmt done"

# Install all packages.
install:
	for i in $(PACKAGES) ; do go install $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): install done"

# Run go generate for all packages.
generate:
	for i in $(PACKAGES) ; do go generate $(GITHUB_LIVEPPROF)/$$i || exit ; done
	@echo "$$(date): generate done"

# Run go meta linter on all packages.
# Returns true whatsoever, read the output to figure out what's wrong.
lint:
	for i in $(PACKAGES) ; do gometalinter --config=.gometalinter.json $$i || exit ; done
	@echo "$$(date): lint done"
