fmt: ;@go fmt ./...

compile: ;@GO111MODULE=auto ;mkdir -p bin/ ; go build  -o ./bin/

build: clean  compile bind

clean: ;@rm -rf ./bin/MailSenitizer

bind: ;@sudo setcap CAP_NET_BIND_SERVICE=+eip ./bin/MailSenitizer