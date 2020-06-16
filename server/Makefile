PREFIX = /usr/local
NAME = sest-server

all:
	go build -o ${NAME}

clean:
	go mod tidy
	go clean

install:
	cp ./${NAME} ../bin/${NAME}
	cp ./${NAME} ${PREFIX}/bin/${NAME}
	chmod 755 ${PREFIX}/bin/${NAME}

uninstall:
	rm ${PREFIX}/bin/${NAME}
