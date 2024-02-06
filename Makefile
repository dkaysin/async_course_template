build:
	go1.21.6 build -o ./bin/server .

run: build
	source .env && ./bin/server
