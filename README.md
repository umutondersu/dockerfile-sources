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

- Go 1.x or higher
- GitHub access token (optional, for higher API rate limits)

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
3. Aggregates results into a structured JSON output
4. Prints final JSON to stdout

- Comprehensive error checking at each step
- Validation of input data format
- Graceful failure handling for GitHub API rate limits and network errors
- Proper error propagation through the application
- Smart retry mechanism with exponential backoff for transient failures
- Detailed error types for better error handling and debugging:
  - Rate limit errors with reset time information
  - GitHub API errors with status codes and messages
  - Network errors with automatic retries
- Maximum retry duration of 30 seconds with increasing intervals
