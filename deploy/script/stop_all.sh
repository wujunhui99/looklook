#!/bin/bash
# deploy/script/stop_all.sh

# Function to stop a specific service
stop_service() {
    local name=$1
    local pattern=$2
    
    echo "Stopping $name..."
    # Kill the go run process
    pkill -f "go run $pattern"
    # Kill the run.sh process
    pkill -f "run.sh $name"
    # Kill legacy make process if any
    pkill -f "make run SERVICE=$name"
}

# Stop all services
stop_all() {
    echo "Stopping all services..."

    # RPCs
    stop_service "usercenter.rpc" "app/usercenter/cmd/rpc/usercenter.go"
    stop_service "travel.rpc" "app/travel/cmd/rpc/travel.go"
    stop_service "payment.rpc" "app/payment/cmd/rpc/payment.go"
    stop_service "order.rpc" "app/order/cmd/rpc/order.go"

    # APIs
    stop_service "usercenter-api" "app/usercenter/cmd/api/usercenter.go"
    stop_service "travel-api" "app/travel/cmd/api/travel.go"
    stop_service "payment-api" "app/payment/cmd/api/payment.go"
    stop_service "order-api" "app/order/cmd/api/order.go"

    # Consumers/Jobs
    stop_service "order-mq" "app/order/cmd/mq/order.go"
    stop_service "mqueue-job" "app/mqueue/cmd/job/mqueue.go"
    stop_service "mqueue-scheduler" "app/mqueue/cmd/scheduler/mqueue.go"

    echo "All services stopped."
}

# Check if arguments are provided
if [ $# -eq 0 ]; then
    stop_all
else
    # Allow stopping specific services passed as arguments
    for service in "$@"; do
        case $service in
            "usercenter.rpc") stop_service "usercenter.rpc" "app/usercenter/cmd/rpc/usercenter.go" ;;
            "travel.rpc") stop_service "travel.rpc" "app/travel/cmd/rpc/travel.go" ;;
            "payment.rpc") stop_service "payment.rpc" "app/payment/cmd/rpc/payment.go" ;;
            "order.rpc") stop_service "order.rpc" "app/order/cmd/rpc/order.go" ;;
            "usercenter-api") stop_service "usercenter-api" "app/usercenter/cmd/api/usercenter.go" ;;
            "travel-api") stop_service "travel-api" "app/travel/cmd/api/travel.go" ;;
            "payment-api") stop_service "payment-api" "app/payment/cmd/api/payment.go" ;;
            "order-api") stop_service "order-api" "app/order/cmd/api/order.go" ;;
            "order-mq") stop_service "order-mq" "app/order/cmd/mq/order.go" ;;
            "mqueue-job") stop_service "mqueue-job" "app/mqueue/cmd/job/mqueue.go" ;;
            "mqueue-scheduler") stop_service "mqueue-scheduler" "app/mqueue/cmd/scheduler/mqueue.go" ;;
            *) echo "Unknown service: $service" ;;
        esac
    done
fi
