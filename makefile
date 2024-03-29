.PHONY: test
test:
	go test -v -run TestGetStockQuote
	go test -v -run TestGetStockQuotes

.PHONY: bingo_rpi
bingo_rpi:
	go mod vendor
	GOOS=linux GOARCH=arm64 go build -o bingo_rpi .

.PHONY: bingo
bingo:
	go mod vendor
	go build -o bingo .

.PHONY: clean
clean:
	rm -rf bingo_rpi bingo_single bingo