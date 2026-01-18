#!/bin/bash
# deploy/script/list_services.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(dirname "$(dirname "$SCRIPT_DIR")")
cd "$PROJECT_ROOT"

echo "Available microservices:"
echo ""
for yaml in app/*/cmd/*/etc/*.yaml; do
    if [ -f "$yaml" ]; then
        service_name=$(grep "^Name:" "$yaml" | awk '{print $2}')
        service_dir=$(dirname "$yaml" | sed 's|/etc$||')
        app_name=$(echo "$service_dir" | awk -F/ '{print $2}')
        service_type=$(echo "$service_dir" | awk -F/ '{print $4}')
        echo "  - $service_name (app/$app_name/cmd/$service_type)"
    fi
done
echo ""
echo "Usage: make run SERVICE=<service-name> [ENV=local|docker|k8s]"
