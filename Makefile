GOCTL_HOME=$(shell pwd)/deploy/goctl/.goctl/1.8.3
CACHE ?= true

# Generate model code for all microservices
.PHONY: gen-model-all
gen-model-all:
	@./deploy/script/gen_model.sh all $(CACHE)

# Generate model code for specified microservice
.PHONY: gen-model
gen-model:
	@./deploy/script/gen_model.sh $(SERVICE) $(CACHE)

# Run microservice
ENV ?= local
.PHONY: run
run:
	@./deploy/script/run.sh $(SERVICE) $(ENV)

# List all available microservices
.PHONY: list-services
list-services:
	@./deploy/script/list_services.sh

# Help information
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Model Generation:"
	@echo "  make gen-model-all              - Generate model code for all microservices"
	@echo "  make gen-model SERVICE=<name>   - Generate model code for specified microservice"
	@echo ""
	@echo "Run Microservices:"
	@echo "  make run SERVICE=<name>         - Run specified microservice"
	@echo "  make list-services              - List all available microservices"