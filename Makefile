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

# Help information
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make gen-model-all              - Generate model code for all microservices (cache enabled by default)"
	@echo "  make gen-model-all CACHE=false  - Generate model code for all microservices (cache disabled)"
	@echo "  make gen-model SERVICE=<name>   - Generate model code for specified microservice (cache enabled by default)"
	@echo "  make gen-model SERVICE=<name> CACHE=false - Generate model code for specified microservice (cache disabled)"
	@echo ""
	@echo "Examples:"
	@echo "  make gen-model SERVICE=usercenter"
	@echo "  make gen-model SERVICE=order CACHE=false"
	@echo "  make gen-model-all"
