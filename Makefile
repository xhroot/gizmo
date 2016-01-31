all: test

deps:
	go get -d -v github.com/xhroot/gizmo/...

updatedeps:
	go get -d -v -u -f github.com/xhroot/gizmo/...

testdeps:
	go get -d -v -t github.com/xhroot/gizmo/...

updatetestdeps:
	go get -d -v -t -u -f github.com/xhroot/gizmo/...

build: deps
	go build github.com/xhroot/gizmo/...

install: deps
	go install github.com/xhroot/gizmo/...

lint: testdeps
	go get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v '\.pb\.go\|\.pb\.gw\.go\|examples\|pubsub\/awssub_test\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

vet: testdeps
	go vet github.com/xhroot/gizmo/...

errcheck: testdeps
	go get -v github.com/kisielk/errcheck
	errcheck -ignoretests github.com/xhroot/gizmo/...

pretest: lint vet # errcheck

test: testdeps pretest
	go test github.com/xhroot/gizmo/...

clean:
	go clean -i github.com/xhroot/gizmo/...

coverage: testdeps
	./coverage.sh --coveralls

.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	clean \
	coverage
