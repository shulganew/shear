# Shortener with Yandex Practicum
## X-Real-IP with curl
```
curl --header "X-Real-IP: 192.168.2.2" localhost:8080/api/internal/stats
```
## Create sertificates

Generate private key:
```
openssl genrsa -out server.pem 2048
```
Generate CSR: (In the "Common Name" set the domain of your service provider app)
```
openssl req -new -key server.pem -out server.csr
```

Generate Self Signed Cert
```
openssl x509 -req -days 365 -in server.csr -signkey server.pem -out server.crt
```

## LD flags
```
-ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=LTS version"
go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=LTS version'" ./cmd/shortener/... 
```
## Custom Analyzers
```
https://golangci-lint.run/usage/linters/
https://github.com/kyoh86/exportloopref
https://github.com/leighmcculloch/gochecknoglobals
https://github.com/butuzov/ireturn
```

Build and use
```
go build -o ./cmd/staticlint/mycheck ./cmd/staticlint/...
./cmd/staticlint/mycheck ./cmd/... ./internal/...
```
## Swagger

```go
go install github.com/swaggo/swag/cmd/swag@latest
swag init --output ./swagger/
swag init --output ./swagger/ --parseDependency --parseInternal  -g cmd/shortener/main.go
```
Add import to generated file for web API!!!
```go
_ "github.com/shulganew/shear.git/docs"
```
Details swagger install:

[See web swagger install to project: ](https://github.com/swaggo/http-swagger)


## Godoc

```
go install -v golang.org/x/tools/cmd/godoc@latest
#####godoc#####
export PATH="$GOPATH/bin:$PATH"
godoc -http=:8085
```
Show internal packages
http://localhost:8085/pkg/?m=all

## Profiling
pprof is a tool for visualization and analysis of profiling data
### Install
```go
go install github.com/google/pprof@latest
```

### WEB
```go
http://127.0.0.1:8080/debug/pprof/
go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/profile 
```

### Consloe
```go
go tool pprof -seconds=30 http://localhost:8080/debug/pprof/profile
```
```cmd
list foo
top foo
```

### Benchmark shortener API with Vegeta
https://github.com/tsenart/vegeta
[Vegeta project: ](https://github.com/tsenart/vegeta)
```
echo "GET http://localhost:8080" | vegeta attack -duration=5s -rate=5 | vegeta report --type=text
```
### Memory (use URL)
```
go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/heap
```

### Memory (use local file)
```
curl -sK -v http://localhost:8080/debug/pprof/heap > profiles/base.pprof
go tool pprof -http=":9090" -seconds=30 profiles/base.pprof 
curl -sK -v http://localhost:8080/debug/pprof/heap > profiles/result.pprof
go tool pprof -http=":9090" -seconds=30 profiles/result.pprof 

pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
pprof -top profiles/result.pprof
```
Save as image:
```
go tool pprof -png profiles/result.pprof > profiles/result.png
go tool pprof -svg profiles/result.pprof > profiles/result.svg
```

## benchmark
```
go test -bench  . ./internal/service/
go test -bench=. -benchmem ./internal/service/
go test -bench=. -cpuprofile ./internal/service/
```
Профиль использования процессорного времени.
```
go test -bench=. -cpuprofile profiles/cpu.out ./internal/service/
```
```
$ go tool pprof ./profiles/cpu.out
File: service.test
Type: cpu
Time: Mar 5, 2024 at 10:44am (MSK)
Duration: 12.23s, Total samples = 14.15s (115.70%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
(pprof)
```
Профиль использования памяти.
```

go test -bench=. -memprofile profiles/memory.out ./internal/service/
$ go tool pprof ./profiles/memory.out
File: service.test
Type: cpu
Time: Mar 5, 2024 at 10:44am (MSK)
Duration: 12.23s, Total samples = 14.15s (115.70%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
(pprof)
```

## cmd commands for test purposes
### GET and POST handles
```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/ -d "https://yandex1.ru"
curl -v -H "Content-Type: text/plain" http://localhost:8080/hjnFtibr

curl -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"https://practicum1.yandex1.ru"}'
curl -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"http://liceih591s.com/rmqtluduv3fe8t/qtefpaham0"}'
```
### gzip
```bash
add --compressed key, this include accept encoding header
curl --compressed -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"https://practicum1.yandex1.ru"}' | gunzip

//send gzip body to server
echo '{"url":"https://practicum1.yandex1.ru"}' | gzip > body.gz
curl --compressed -v -X POST http://localhost:8080/api/shorten -H'Content-Encoding: gzip' --data-binary @body.gz | gunzip


set SERVER_ADDRESS=localhost:8080
echo %SERVER_ADDRESS%
```
# Git
```
git push -u origin iter5
git switch iter1

git add .
git commit -m "refactor del handler"
git push -u origin iter15

git log
diff 80846d62286fd8d87e16f9ae833f3e859ab8ecaf 0b2e0a457317162be17339da7c557f3d56c3db8a
```
# UUID
```
https://pkg.go.dev/github.com/google/uuid#section-readme
```

# PostGreSQL

```bash
Connection string:
postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
postgresql://short:1@localhost/short
~/.bashrc
export DATABASE_DSN=postgresql://short:1@localhost/short
```

```bash
# Backup
~/.bash_profile
export FILE_STORAGE_PATH=/tmp/short-url-db.json
```

# Tests

## Run static test localy

```bash
go vet -vettool=$(which statictest) ./...
shortenertestbeta -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=.
```

## Iterantion tests
go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go

## my links to useful sites

## Test cover
```
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov
```
```
go test -v -coverpkg=./... -coverprofile=profile.cov ./...;cat profile.cov | grep -v "mock\|vegeta" > nomocksmigrations.cov;go tool cover -func nomocksmigrations.cov
```

# Use autotest local 
https://github.com/nektos/act


# My links
```bash
https://github.com/golang/go/wiki/CodeReviewComments#receiver-type

status code during the test issue:
https://github.com/gin-gonic/gin/issues/1120

hadle func with interface
https://ru.hexlet.io/courses/go-web-development/lessons/local-persistence/theory_unit

Go Interface in detail:
https://research.swtch.com/interfaces

file lseek
https://www.opennet.ru/docs/RUS/zlp/005.html
```

# Lint install
```bash
go install golang.org/x/tools/gopls@latest
//========================================
run ci lint

sudo snap install golangci-lint
golangci-lint run
golangci-lint --help
golangci-lint run -v
golangci-lint run

```

## Mock generate 

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

```bash
mockgen -source=internal/service/shortener.go \
    -destination=internal/service/mocks/shortener_mock.gen.go \
    -package=mocks
```

## Запуск Postgres в контейнере

Для запуска и остановки Postgres в контейнере выполнятьются скрипты создания и миграции базы в make-файле:
* Инициализация
```bash
make pg
```
* Миграция goose
```bash
https://github.com/pressly/goose
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="postgresql://market:1@localhost/market"
GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://market:1@localhost/market" goose up
GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://postgres:postgres@postgres/praktikum" goose -dir ./migrations  up
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgresql://postgres:postgres@postgres/praktikum"
```

* Остановка и удаление контейнера
```bash
make pg-stop
```

# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.


## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```bash
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```bash
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


