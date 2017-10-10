install:
	cp example.cnf /etc/caserver.cnf && mkdir -p /var/lib/caserver/ && cp build/caserver /usr/local/bin/caserver
build:
	go build -o build/caserver -v -x -ldflags '-w -s -linkmode external -extldflags "-static"' main.go && strip build/caserver