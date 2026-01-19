#!/bin/bash

# deploy/script/gen_model.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(dirname "$(dirname "$SCRIPT_DIR")")
cd "$PROJECT_ROOT"

GOCTL_HOME="$PROJECT_ROOT/deploy/goctl/.goctl/1.8.3"
SERVICE=$1
CACHE=${2:-true}

if [ "$CACHE" = "true" ]; then
    CACHE_FLAG="-c"
else
    CACHE_FLAG=""
fi

if [ -z "$SERVICE" ]; then
     echo "Error: Service name or 'all' required."
     echo "Usage: ./deploy/script/gen_model.sh <service_name|all> [cache_enabled:true|false]"
     exit 1
fi

if [ "$SERVICE" = "all" ]; then
    echo "Starting to generate model code for all microservices..."
    for sql_file in deploy/sql/*.sql; do
        if [ -f "$sql_file" ]; then
            filename=$(basename "$sql_file" .sql)
            # filename format expected: prefix_service (e.g., looklook_order)
            service=$(echo $filename | awk -F'_' '{print $NF}')
            echo "Generating model code for $service..."
            mkdir -p "app/$service/model"
            goctl model mysql ddl -src="$sql_file" -dir="app/$service/model" --home="$GOCTL_HOME" $CACHE_FLAG
        fi
    done
    echo "All model code generation completed!"
else
    # Single service
    sql_file=$(ls deploy/sql/*_${SERVICE}.sql 2>/dev/null | head -n 1)
    
    if [ -z "$sql_file" ]; then
        echo "Error: SQL file not found for microservice $SERVICE"
        exit 1
    fi
    
    echo "Generating model code for $SERVICE..."
    mkdir -p "app/$SERVICE/model"
    goctl model mysql ddl -src="$sql_file" -dir="app/$SERVICE/model" --home="$GOCTL_HOME" $CACHE_FLAG
    echo "Model code generation for $SERVICE completed!"
fi
