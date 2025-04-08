#!/bin/bash

IMAGE_NAME=dockerfile-sources

# Start minikube if not already running
if ! minikube status >/dev/null 2>&1; then
    echo "Starting minikube..."
    minikube start
fi

# Set docker env to use minikube's docker daemon
eval "$(minikube -p minikube docker-env --shell=bash)"
#NOTE: use --unset to revert back to local docker daemon

# Check if image exists
if ! docker images "$IMAGE_NAME":latest | grep "$IMAGE_NAME" >/dev/null 2>&1; then
    echo "Building dockerfile-sources:latest image..."
    docker build -t dockerfile-sources:latest .
else
    echo "Image dockerfile-sources:latest already exists"
fi

# Apply Kubernetes resources
echo "Applying Kubernetes resources..."

if ! kubectl apply -f k8s/secret.yaml; then
    echo "Failed to apply secret configuration"
fi

if ! kubectl apply -f k8s/job.yaml; then
    echo "Failed to apply job configuration"
    exit 1
fi

echo "Job started successfully"
