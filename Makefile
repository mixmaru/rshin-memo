build_webui:
	cd ./web/rshinmemo && go build

run_webui: build_webui
	cd ./web/rshinmemo && ./rshinmemo

run_webui_with_test_data: build_webui
	cd ./web/rshinmemo && ./rshinmemo -datadir=./testdata

build_gocui:
	cd ./cmd/rshinmemo && go build

exec_gocui: build_gocui
	cd ./cmd/rshinmemo && ./rshinmemo

build_tview:
	cd ./tui/rshinmemo && go build

exec_tview: build_tview
	cd ./tui/rshinmemo && ./rshinmemo

clean:
	rm -f ./cmd/rshinmemo/rshinmemo
	rm -f ./tui/rshinmemo/rshinmemo

test:
	go test ./... -count=1 -cover
