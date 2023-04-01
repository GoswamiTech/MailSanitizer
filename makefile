fmt: ;@go fmt ./...

build: ;@GO111MODULE=on ;mkdir -p bin/ ; go build  -o ./bin/

clean: ;@rm -rf ./bin/MailSenitizer