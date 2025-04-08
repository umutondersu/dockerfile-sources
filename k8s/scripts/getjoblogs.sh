#!/bin/bash

JOB_NAME=dockerfile-sources-job

# Wait for the pod to be created and get its name
POD_NAME=$(kubectl get pods --selector=job-name="$JOB_NAME" --output=jsonpath='{.items[*].metadata.name}')

# Wait until the pod exists
while [ "$POD_NAME" = "" ]; do
    echo "Waiting for pod to be created..."
    sleep 2
    POD_NAME=$(kubectl get pods --selector=job-name="$JOB_NAME" --output=jsonpath='{.items[*].metadata.name}')
done

echo "Found pod: $POD_NAME"

# Follow the logs
kubectl logs -f "$POD_NAME"

