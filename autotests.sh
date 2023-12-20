#!/bin/bash
export PATH="/home/igor/Desktop/code/go-autotests-0.10.2/bin:$PATH"
unset DATABASE_DSN
rm -f /tmp/short-url-db.json
function check(){
	res=""
	if [ $2 -ne 0 ]; then res=$(echo "$1: {$res} Error! $2"); echo "ERROR!  Iter:" $res;exit 1; else res=$(echo "$1: ${res} PASS "); fi
	echo "Iter:" $res
}

go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go
go vet -vettool=$(which statictest) ./...
check S $? 
TEMP_FILE=$(random tempfile)
#1
shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
check 1 $?
#2
shortenertestbeta -test.v -test.run=^TestIteration2$ -source-path=. > /dev/null
check 2 $?
#3
shortenertestbeta -test.v -test.run=^TestIteration3$ -source-path=. > /dev/null
check 3 $?
#4
SERVER_PORT=$(random unused-port)
shortenertestbeta -test.v -test.run=^TestIteration4$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT > /dev/null
check 4 $?
#5
SERVER_PORT=$(random unused-port)
shortenertestbeta -test.v -test.run=^TestIteration5$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT > /dev/null
check 5 $?
#6
shortenertestbeta -test.v -test.run=^TestIteration6$ -source-path=. > /dev/null
check 6 $?
#7
shortenertestbeta -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=. > /dev/null
check 7 $?
#8
shortenertestbeta -test.v -test.run=^TestIteration8$ -binary-path=cmd/shortener/shortener \
> /dev/null
check 8 $?
#9
shortenertestbeta -test.v -test.run=^TestIteration9$  -binary-path=cmd/shortener/shortener -source-path=. -file-storage-path=/tmp/short-url-db.json > /dev/null
check 9 $?
#10
#
          shortenertestbeta -test.v -test.run=^TestIteration10$ \
              -binary-path=cmd/shortener/shortener \
              -source-path=. \
              -database-dsn='postgres://short:1@postgres:5432/praktikum?sslmode=disable'
check 10 $?
#11
          shortenertestbeta -test.v -test.run=^TestIteration11$ \
              -binary-path=cmd/shortener/shortener \
              -database-dsn='postgres://short:1@postgres:5432/praktikum?sslmode=disable'
check 11 $?
#12
          shortenertestbeta -test.v -test.run=^TestIteration12$ \
              -binary-path=cmd/shortener/shortener \
              -database-dsn='postgres://short:1@postgres:5432/praktikum?sslmode=disable'
check 12 $?
#13
          shortenertestbeta -test.v -test.run=^TestIteration13$ \
              -binary-path=cmd/shortener/shortener \
              -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
check 13 $?
#14
#          shortenertestbeta -test.v -test.run=^TestIteration14$ \
#              -binary-path=cmd/shortener/shortener \
#              -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
#15
#          shortenertestbeta -test.v -test.run=^TestIteration15$ \
#              -binary-path=cmd/shortener/shortener \
#              -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
#16
#shortenertestbeta -test.v -test.run=^TestIteration16$ -source-path=.
#17
#shortenertestbeta -test.v -test.run=^TestIteration17$ -source-path=.
#18
#         shortenertestbeta -test.v -test.run=^TestIteration18$ \
#             -source-path=. \
