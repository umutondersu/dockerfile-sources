#!/bin/bash

JOB_NAME=dockerfile-sources-job

# Delete the Kubernetes job and secret
echo "Cleaning up Kubernetes resources..."
kubectl delete job "$JOB_NAME"
kubectl delete secret github-token

echo "Job and related resources cleaned up successfully"
