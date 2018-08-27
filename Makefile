.PHONY: run
.SILENT: run test eqso

test:
	go test -run TestParser

run: eqso
	./eqso

eqso: eqso.go parser.go lexer.go
	go build -i -o eqso
