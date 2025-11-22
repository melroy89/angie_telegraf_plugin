binary = angie_telegraf

all: clean build

build:
	go build -o $(binary) cmd/main.go

clean:
	rm -rf ./$(binary)

restart: 
	sudo systemctl restart telegraf

rundev: build
	./angie_telegraf -config ./dev.conf

