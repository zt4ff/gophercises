package urlshort

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func buildMap(pathUrls []pathUrl, pathsToUrls map[string]string) {
	for _, x := range pathUrls {
		pathsToUrls[x.Path] = x.URL
	}
}

func ParseYAML(filepath string, pathsToUrls *map[string]string) error {
	var pathUrls []pathUrl
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &pathUrls)
	if err != nil {
		return err
	}

	buildMap(pathUrls, *pathsToUrls)
	return nil
}

func ParseJSON(filepath string, pathsToUrls *map[string]string) error {
	var pathUrls []pathUrl
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &pathUrls)
	if err != nil {
		return err
	}

	buildMap(pathUrls, *pathsToUrls)
	return nil
}
