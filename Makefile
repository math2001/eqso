.PHONY: run
.SILENT: run test eqso

test:
	go test

run: eqso
	./eqso

eqso: eqso.go
	go build -i
