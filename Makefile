include project.conf

build:
	go build -mod vendor -o _build/moekyuniversity cmd/moekyuniversity.go
	cp -r static _build/static

start:
	./_build/moekyuniversity -static-dir _build/static

clean:
	rm -rf _build

test:
	go test -mod vendor -v ./...

generate-test-data:
	go run -mod vendor cmd/generatetestdata.go -cards-file data/cards.json