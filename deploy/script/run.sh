#!/bin/bash

# deploy/script/run.sh

# Resolve project root
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(dirname "$(dirname "$SCRIPT_DIR")")

cd "$PROJECT_ROOT"

SERVICE=$1
ENV=${2:-local}
ENV_FILE=".env.$ENV"

if [ -z "$SERVICE" ]; then
    echo "Error: Please specify SERVICE parameter"
    echo "Usage: ./deploy/script/run.sh <service-name> [env]"
    exit 1
fi

if [ "$ENV" != "local" ] && [ "$ENV" != "docker" ] && [ "$ENV" != "k8s" ]; then
    echo "Error: ENV must be one of: local, docker, k8s"
    exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: Environment file $ENV_FILE not found"
    exit 1
fi

echo "Running microservice: $SERVICE in $ENV environment"
echo "Environment file: $ENV_FILE"

service_found=false

# Search for the service in yaml files
for yaml in app/*/cmd/*/etc/*.yaml; do
    if [ -f "$yaml" ]; then
        # Extract Name from yaml
        service_name=$(grep "^Name:" "$yaml" | awk '{print $2}')
        
        if [ "$service_name" = "$SERVICE" ]; then
            service_found=true
            
            # Determine paths based on yaml location: app/<app_name>/cmd/<service_type>/etc/<file>
            app_name=$(echo "$yaml" | awk -F/ '{print $2}')
            service_type=$(echo "$yaml" | awk -F/ '{print $4}')
            
            main_file="app/$app_name/cmd/$service_type/$app_name.go"
            
            if [ ! -f "$main_file" ]; then
                echo "Error: Main file $main_file not found"
                exit 1
            fi
            
            echo "Loading environment variables from $ENV_FILE..."
            echo "Starting $main_file with config $yaml..."
            
            # Export environment variables and run
            # We use eval to handle the export correctly with xargs output
            export $(cat $ENV_FILE | grep -v '^#' | grep -v '^$' | xargs)
            go run "$main_file" -f "$yaml"
            break
        fi
    fi
done

if [ "$service_found" = "false" ]; then
    echo "Error: Service $SERVICE not found"
    exit 1
fi
