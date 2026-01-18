#!/bin/bash

# deploy/script/start_all.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# Function to run a service
run_service() {
    service=$1
    echo "Starting $service..."
    "$SCRIPT_DIR/run.sh" "$service" &
}

# Start RPCs
run_service "usercenter.rpc"
run_service "travel.rpc"
run_service "payment.rpc"
run_service "order.rpc"
sleep 5

# Start APIs and Consumers
run_service "usercenter-api"
run_service "travel-api"
run_service "payment-api"
run_service "order-api"
run_service "order-mq"
run_service "mqueue-job"
run_service "mqueue-scheduler"

echo "All services started in background."
