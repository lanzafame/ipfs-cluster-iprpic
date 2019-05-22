install: build
	scp ./iprpid pi@$(RPI):/home/pi/

build:
	env GOOS=linux GOARCH=arm GO111MODULE=on go build ./cmd/iprpid

clean:
	rm ./iprpid
