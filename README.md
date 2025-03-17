# Dockerfile Sources

A Go application that extracts and lists all container images used in Dockerfiles across GitHub repositories.

## Overview

This tool scans Dockerfiles from specified GitHub repositories and generates a JSON output that maps repositories to their container image dependencies.

## Features

- Authenticates with GitHub (optional for higher API limits)
- Scans Dockerfiles in the specified repositories
- Extracts all container images (FROM statements) from Dockerfiles
- Generates JSON output mapping repositories to their container images

## Example Output

```json
{
  "data": {
    "owner/repo1": {
      "Dockerfile": ["nginx:1.19", "node:14"],
      "services/Dockerfile": ["python:3.9"]
    },
    "owner/repo2": {
      "Dockerfile.prod": ["golang:1.16"]
    }
  }
}
```

## Prerequisites

- GitHub access token with repo scope (optional, for private repository access and higher API rate limits)

## Environment Variables

The application requires the following environment variables:

- `REPOSITORY_LIST_URL`: URL to the text file containing repository sources
- `GITHUB_ACCESS_TOKEN`: GitHub access token for API authentication (optional)

## Installation

```bash
git clone https://github.com/umutondersu/dockerfile-sources
cd dockerfile-sources
go build ./cmd/dockerfile-sources
# dockerfile-sources
```

## Kubernetes Deployment

The application can be deployed as a Kubernetes job. The necessary configuration files and helper scripts are provided in the `k8s` directory.

### Prerequisites

- Minikube
- kubectl
- Docker

### Configuration Files

- `k8s/job.yaml`: Kubernetes job configuration with resource limits and timeouts
- `k8s/secret.yaml`: Template for using `GITHUB_ACCESS_TOKEN`
- `k8s/scripts/`: Helper scripts for managing the job

### Helper Scripts

The following scripts are provided to simplify deployment and monitoring:

- `startjob.sh`: Handles Minikube startup, builds Docker image, and deploys the job
- `getjoblogs.sh`: Monitors job logs by automatically finding the pod
- `endjob.sh`: Cleans up job resources

### Running the Job

1. Make the scripts executable:

```bash
chmod +x k8s/scripts/*.sh
```

2. Start the job:

```bash
./k8s/scripts/startjob.sh
```

3. Monitor logs (in a separate terminal):

```bash
./k8s/scripts/getjoblogs.sh
```

4. Clean up when done:

```bash
./k8s/scripts/endjob.sh
```

## Implementation Details

### Package Structure

The application is organized into several packages, each with a specific responsibility:

#### `internal/input`

- Handles parsing of repository source data
- Contains the `Source` struct that represents a GitHub repository with owner, repo name, and commit SHA
- Implements HTTP response fetching for repository list URL
- Uses regex pattern matching to extract repository information

#### `internal/ghdocker`

- Manages GitHub API interactions through a custom client
- Responsible for scanning repositories for Dockerfile presence
- Extracts container image information from Dockerfiles
- The `DockerFile` struct maintains the relationship between sources and their container images
- Implements concurrent processing of multiple repositories

#### `internal/jsonoutput`

- Handles the conversion of Dockerfile data into the final JSON format
- Implements the `OutputData` structure that maps repositories to their Dockerfile paths and images
- Provides JSON generation functionality with proper error handling
- Ensures consistent output format as shown in the example above

#### `cmd/dockerfile-sources`

- Contains the main application entry point
- Orchestrates the workflow:
  1. Reads environment variables
  2. Fetches repository list
  3. Initializes GitHub client
  4. Processes Dockerfiles
  5. Generates and outputs JSON

## Technical Design

### Data Flow

1. Application reads repository list from provided URL
2. For each repository:
   - Scans for Dockerfile presence
   - Extracts FROM statements to identify container images
   - Uses goroutines and channels for efficient file content retrieval
   - Implements wait groups to ensure all processing completes
3. Aggregates results into a structured JSON output
4. Prints final JSON to stdout

### Performance Features

- Concurrent processing of Dockerfile content
- Channel-based communication for efficient data transfer
- Wait group synchronization for parallel operations
- Optimized memory usage through streaming processing

### Other Features

- Comprehensive error checking at each step
- Validation of input data format
- Operation timeout control (default: 5 minutes)
- Graceful failure handling for GitHub API rate limits and network errors
- Proper error propagation through the application
- Robust error handling system with:
  - Centralized error processing through `handleGitHubResponseError`
  - Smart retry mechanism with exponential backoff for transient failures
  - Standardized error patterns across all GitHub API calls
  - Detailed error types for better debugging:
  - Maximum retry duration of 30 seconds with 100ms initial interval
  - Automatic backoff for retryable errors (500+ status codes)
  - Concurrent operation error capture and reporting
- Consistent exit codes for different error scenarios
