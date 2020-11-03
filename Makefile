default: build

.PHONEY: build
build:
	go build cmd/gocache.go
	go build cmd/goctl.go


test:
	echo test here.