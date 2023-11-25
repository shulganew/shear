#!/bin/bash
go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go
go vet -vettool=$(which statictest) ./...
#shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
#shortenertestbeta -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=.
TEMP_FILE=$(random tempfile)
shortenertestbeta -test.v -test.run=^TestIteration9$  -binary-path=cmd/shortener/shortener -source-path=. -file-storage-path=$TEMP_FILE