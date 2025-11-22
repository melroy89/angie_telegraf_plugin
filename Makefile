bin_dir = ./bin
binary = angie_telegraf_plugin

all: clean build

build:
	mkdir -p $(bin_dir)
	go build -o $(bin_dir)/$(binary) cmd/main.go

clean:
	rm -rf $(bin_dir)

restart: 
	sudo systemctl restart telegraf

rundev: build
	./bin/angie_telegraf_plugin -config ./dev.conf

