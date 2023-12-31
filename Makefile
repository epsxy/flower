.PHONY: release-darwin-amd64, release-windows-amd64

run:
	go run main.go

test:
	bazel test --test_output=errors //...

cov:
	bazel coverage --instrument_test_targets --experimental_cc_coverage --combined_report=lcov //...
# genhtml --output genhtml $(bazel info output_path)/_coverage/_coverage_report.dat

rm-cov:
	rm -rf genhtml

clean:
	rm -rf bin/*.plantuml
	rm -rf genhtml

docker-build:
	docker build . --tag flower:alpha

docker-run:
	docker run -it --rm -v $(CURDIR)/bin:/app/bin flower:alpha /bin/bash

build:
	go build -o ./bin/flower main.go

install:
	go build -o ./bin/flower main.go
	sudo cp ./bin/flower /usr/local/bin/flower

uninstall:
	sudo rm -rf /usr/local/bin/flower

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64 main.go &

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64 main.go &

build-freebsd-amd64:
	GOOS=freebsd GOARCH=amd64 go build -o bin/freebsd-amd64 main.go &

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64 main.go &

release-darwin-amd64: build-darwin-amd64
	mkdir -p release/darwin-amd64
	cp release/Makefile release/darwin-amd64
	cp release/README.md release/darwin-amd64
	cp bin/darwin-amd64 release/darwin-amd64/flower
	tar -czvf release/darwin-amd64.tar.gz release/darwin-amd64
	rm -rf release/darwin-amd64