.PHONY: test
test:
	go test -v -run TestAnimeGAN

.PHONY: bingo_rpi
bingo_rpi:
	go mod vendor
	GOOS=linux GOARCH=arm64 go build -o bingo_rpi .

.PHONY: clean
clean:
	rm -rf bingo_rpi bingo_single bingo