build:
	cd ./cmd/rshinmemo && go build

exec: build
	cd ./cmd/rshinmemo && ./rshinmemo

clean:
	rm -f ./cmd/rshinmemo/rshinmemo

test:
	go test ./... -count=1 -cover