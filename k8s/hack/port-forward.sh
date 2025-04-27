#!/bin/bash

set -e

NAMESPACE=numbers-system

echo "🚀 Starting Port-Forwarding for Publishers and Consumer..."

# Forward each Publisher service
kubectl port-forward svc/publisher-1 8081:8081 -n $NAMESPACE &
kubectl port-forward svc/publisher-2 8082:8082 -n $NAMESPACE &
kubectl port-forward svc/publisher-3 8083:8083 -n $NAMESPACE &
kubectl port-forward svc/publisher-4 8084:8084 -n $NAMESPACE &
kubectl port-forward svc/publisher-5 8085:8085 -n $NAMESPACE &

# Forward Consumer service
kubectl port-forward svc/consumer 9090:9090 -n $NAMESPACE &

echo "Port-Forwards started."
echo "Access Publishers: http://localhost:8081, 8082, 8083, 8084, 8085"
echo "Access Consumer:   http://localhost:9090"
echo ""
echo "Press Ctrl+C to stop forwarding manually."
