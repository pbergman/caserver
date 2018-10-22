install:
	cp example.cnf /etc/caserver.cnf && mkdir -p /var/lib/caserver/ && cp build/caserver /usr/local/bin/caserver
build:
	go build -o build/caserver -v -x -ldflags '-w -s -linkmode external -extldflags "-static"' main.go && strip build/caserver
test:
	if [ ! -d "./cover/" ]; then mkdir ./cover; fi
	for i in $(shell find ./ -type f -name "*_test.go" -exec dirname {} \; | uniq | tr -d './'); do \
		go test -coverprofile ./cover/cover.$$i.out ./$$i; \
		go tool cover -html=./cover/cover.$$i.out -o ./cover/cover.$$i.html; \
	done