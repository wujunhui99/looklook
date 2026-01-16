GOCTL_HOME=$(shell pwd)/deploy/goctl/.goctl/1.8.3
CACHE ?= true

ifeq ($(CACHE),true)
	CACHE_FLAG = -c
else
	CACHE_FLAG =
endif

# Generate model code for all microservices
.PHONY: gen-model-all
gen-model-all:
	@echo "Starting to generate model code for all microservices..."
	@for sql_file in deploy/sql/*.sql; do \
		if [ -f "$$sql_file" ]; then \
			filename=$$(basename $$sql_file .sql); \
			service=$$(echo $$filename | awk -F'_' '{print $$NF}'); \
			echo "Generating model code for $$service..."; \
			mkdir -p app/$$service/model; \
			goctl model mysql ddl -src="$$sql_file" -dir="app/$$service/model" --home="$(GOCTL_HOME)" $(CACHE_FLAG); \
		fi \
	done
	@echo "All model code generation completed!"

# Generate model code for specified microservice
# Usage: make gen-model SERVICE=usercenter
# Disable cache: make gen-model SERVICE=usercenter CACHE=false
.PHONY: gen-model
gen-model:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: Please specify SERVICE parameter"; \
		echo "Usage: make gen-model SERVICE=usercenter"; \
		echo "Disable cache: make gen-model SERVICE=usercenter CACHE=false"; \
		exit 1; \
	fi
	@sql_file=$$(ls deploy/sql/*_$(SERVICE).sql 2>/dev/null | head -n 1); \
	if [ -z "$$sql_file" ]; then \
		echo "Error: SQL file not found for microservice $(SERVICE)"; \
		echo "Available microservices:"; \
		ls deploy/sql/*.sql | xargs -n 1 basename | sed 's/.*_/  - /' | sed 's/.sql//'; \
		exit 1; \
	fi; \
	echo "Generating model code for $(SERVICE)..."; \
	mkdir -p app/$(SERVICE)/model; \
	goctl model mysql ddl -src="$$sql_file" -dir="app/$(SERVICE)/model" --home="$(GOCTL_HOME)" $(CACHE_FLAG); \
	echo "Model code generation for $(SERVICE) completed!"

# Run microservice
# ENV: local (default), docker, k8s
# SERVICE: microservice name (e.g., usercenter-api, usercenter.rpc, travel-api, travel-rpc, order.rpc, payment.rpc)
# Usage: make run SERVICE=usercenter-api ENV=local
ENV ?= local
ENV_FILE = .env.$(ENV)

.PHONY: run
run:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: Please specify SERVICE parameter"; \
		echo "Usage: make run SERVICE=<service-name> [ENV=local|docker|k8s]"; \
		echo ""; \
		echo "Available services:"; \
		for yaml in app/*/cmd/*/etc/*.yaml; do \
			if [ -f "$$yaml" ]; then \
				service_name=$$(grep "^Name:" "$$yaml" | awk '{print $$2}'); \
				echo "    - $$service_name"; \
			fi \
		done; \
		exit 1; \
	fi
	@if [ "$(ENV)" != "local" ] && [ "$(ENV)" != "docker" ] && [ "$(ENV)" != "k8s" ]; then \
		echo "Error: ENV must be one of: local, docker, k8s"; \
		exit 1; \
	fi
	@if [ "$(ENV)" != "local" ]; then \
		echo "Error: Only 'local' environment is currently supported"; \
		echo "Docker and k8s environments will be implemented in the future"; \
		exit 1; \
	fi
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "Error: Environment file $(ENV_FILE) not found"; \
		exit 1; \
	fi
	@echo "Running microservice: $(SERVICE) in $(ENV) environment"
	@echo "Environment file: $(ENV_FILE)"
	@service_found=false; \
	for yaml in app/*/cmd/*/etc/*.yaml; do \
		if [ -f "$$yaml" ]; then \
			service_name=$$(grep "^Name:" "$$yaml" | awk '{print $$2}'); \
			if [ "$$service_name" = "$(SERVICE)" ]; then \
				service_found=true; \
				service_dir=$$(dirname "$$yaml" | sed 's|/etc$$||'); \
				app_name=$$(echo "$$service_dir" | awk -F/ '{print $$2}'); \
				service_type=$$(echo "$$service_dir" | awk -F/ '{print $$4}'); \
				main_file="app/$$app_name/cmd/$$service_type/$$app_name.go"; \
				if [ ! -f "$$main_file" ]; then \
					echo "Error: Main file $$main_file not found"; \
					exit 1; \
				fi; \
				echo "Loading environment variables from $(ENV_FILE)..."; \
				echo "Starting $$main_file with config $$yaml..."; \
				export $$(cat $(ENV_FILE) | grep -v '^#' | grep -v '^$$' | xargs) && \
				go run "$$main_file" -f "$$yaml"; \
				break; \
			fi \
		fi \
	done; \
	if [ "$$service_found" = "false" ]; then \
		echo "Error: Service $(SERVICE) not found"; \
		echo "Please use 'make run' without parameters to see available services"; \
		exit 1; \
	fi

# List all available microservices
.PHONY: list-services
list-services:
	@echo "Available microservices:"
	@echo ""
	@for yaml in app/*/cmd/*/etc/*.yaml; do \
		if [ -f "$$yaml" ]; then \
			service_name=$$(grep "^Name:" "$$yaml" | awk '{print $$2}'); \
			service_dir=$$(dirname "$$yaml" | sed 's|/etc$$||'); \
			app_name=$$(echo "$$service_dir" | awk -F/ '{print $$2}'); \
			service_type=$$(echo "$$service_dir" | awk -F/ '{print $$4}'); \
			echo "  - $$service_name (app/$$app_name/cmd/$$service_type)"; \
		fi \
	done
	@echo ""
	@echo "Usage: make run SERVICE=<service-name> [ENV=local|docker|k8s]"
	@echo "Example: make run SERVICE=usercenter-api ENV=local"

# Help information
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Model Generation:"
	@echo "  make gen-model-all              - Generate model code for all microservices (cache enabled by default)"
	@echo "  make gen-model-all CACHE=false  - Generate model code for all microservices (cache disabled)"
	@echo "  make gen-model SERVICE=<name>   - Generate model code for specified microservice (cache enabled by default)"
	@echo "  make gen-model SERVICE=<name> CACHE=false - Generate model code for specified microservice (cache disabled)"
	@echo ""
	@echo "Run Microservices:"
	@echo "  make run SERVICE=<name> [ENV=local|docker|k8s]  - Run specified microservice (default: ENV=local)"
	@echo "  make list-services              - List all available microservices"
	@echo ""
	@echo "Examples:"
	@echo "  make gen-model SERVICE=usercenter"
	@echo "  make gen-model SERVICE=order CACHE=false"
	@echo "  make gen-model-all"
	@echo "  make run SERVICE=usercenter-api ENV=local"
	@echo "  make run SERVICE=travel-rpc"
	@echo "  make list-services"
