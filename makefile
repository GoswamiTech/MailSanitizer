default : build

fmt: ;@go fmt ./...

compile: ;@GO111MODULE=on ;mkdir -p bin/ ; go build  -o ./bin/

build: clean fmt compile bind

clean: ;@rm -rf ./bin/MailSanitizer

bind: ;@sudo setcap CAP_NET_BIND_SERVICE=+eip ./bin/MailSanitizer