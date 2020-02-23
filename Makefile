VERSION:=$(shell /bin/date "+%Y-%m-%d_%H-%M")

run:
	go run -race *.go

build:
	GOARCH=arm GOARM=7 go build -ldflags "-s -w"
	upx smtpd
	go build -ldflags "-s -w -X 'main.VERSION=$(VERSION)'" -o smtpd-x86
	upx smtpd-x86
