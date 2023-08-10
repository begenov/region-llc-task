run:
	docker-compose up

stop:
	docker-compose down

test:
	go test -v -cover -short ./...


.PHONY: test run stop
