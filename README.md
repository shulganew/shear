# cmd commands for test purposes
```
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/ -d "https://yandex1.ru"
curl -v -H "Content-Type: text/plain" http://localhost:8080/hjnFtibr

curl -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"https://practicum1.yandex1.ru"}'
curl -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"http://liceih591s.com/rmqtluduv3fe8t/qtefpaham0"}'


//gzip
//add --compressed key, this include accept encoding header
curl --compressed -v -H "Content-Type: application/json" -X POST http://localhost:8080/api/shorten -d '{"url":"https://practicum1.yandex1.ru"}' | gunzip

//send gzip body to server
echo '{"url":"https://practicum1.yandex1.ru"}' | gzip > body.gz
curl --compressed -v -X POST http://localhost:8080/api/shorten -H'Content-Encoding: gzip' --data-binary @body.gz | gunzip


set SERVER_ADDRESS=localhost:8080
echo %SERVER_ADDRESS%
```
# Git

//git push -u origin iter5
//git checkout -b iter1


# PostGreSQL

```
Connection string:
postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
postgresql://short:1@localhost/short
export DATABASE_DSN=postgresql://short:1@localhost/short
```

# Tests

# Run static test localy

go vet -vettool=$(which statictest) ./...

shortenertestbeta -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=.

# Iterantion tests
go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go

# my links to useful sites

# Use autotest local 
https://github.com/nektos/act


# My links
```
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

# ###############Lint instal##########################
```
go install golang.org/x/tools/gopls@latest
//========================================
run ci lint

sudo snap install golangci-lint
golangci-lint run
golangci-lint --help
golangci-lint run -v
golangci-lint run

```

# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).
