include project.conf

build:
	go build -mod vendor -o _build/moekyuniversity cmd/moekyuniversity.go
	cp -r static _build/static

clean:
	rm -rf _build

start:
	./_build/moekyuniversity -cards-file data/cards.json -backup-dir data/backup -static-dir _build/static

test:
	go test -mod vendor -v ./...

generate-test-data:
	go run -mod vendor cmd/generatetestdata.go -cards-file data/cards.json