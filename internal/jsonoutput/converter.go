package jsonoutput

import (
	"encoding/json"
	"fmt"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
)

// OutputData represents the root JSON structure
type OutputData struct {
	Data map[string]map[string][]string `json:"data"`
}

func MapDockerfilesToOutputData(dockerfiles []ghdocker.DockerFile) (*OutputData, error) {
	output := &OutputData{
		Data: make(map[string]map[string][]string),
	}

	for _, df := range dockerfiles {
		repoKey := fmt.Sprintf("https://github.com/%s/%s.git:%s",
			df.Source.Owner,
			df.Source.Repo,
			df.Source.CommitSha)

		if _, exists := output.Data[repoKey]; !exists {
			output.Data[repoKey] = make(map[string][]string)
		}

		output.Data[repoKey][df.Path] = df.Images
	}

	return output, nil
}

func toJSON(o *OutputData) (string, error) {
	jsonBytes, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func GenerateJSONOutput(dockerfiles []ghdocker.DockerFile) (string, error) {
	output, err := MapDockerfilesToOutputData(dockerfiles)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	jsonStr, err := toJSON(output)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return jsonStr, nil
}
