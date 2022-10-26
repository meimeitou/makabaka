.PHONY: db
db:
	@docker-compose -f examples/docker-compose.yaml up -d

.PHONY: down
down:
	@docker-compose -f examples/docker-compose.yaml down

.PHONY: gen
gen:
	@go run ./gen/main.go 

.PHONY: run
run:
	@go run ./main.go serve examples/config.yaml    
