GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-c-docker: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.55.2 \
        golangci-lint run \
            -c .golangci.yml \

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint 


.PHONY: go-postgres-start
go-postgres-start:
	@if [[ "$(docker images -q shortdb:v1 2> /dev/null)" == "" ]]; then \
		echo "buid image"; \
		#echo '127.0.0.1 postgres' | sudo tee -a /etc/hosts \
		docker build -t shortdb:v1 - <Dockerfile; \
	fi
	docker run -d --name="shortdb" -p 5432:5432 shortdb:v1
	

.PHONY: go-postgres-stop
go-postgres-stop:
	docker stop shortdb
	docker rm shortdb
