package main

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Setting struct {
	Name           string    `yaml:"name"`
	Folder         string    `yaml:"-"`
	Host           string    `yaml:"host"`
	Port           int       `yaml:"port"`
	SwaggerEnabled bool      `yaml:"swaggerEnabled"`
	Requests       []Request `yaml:"requests"`
	Headers        []Header  `yaml:"headers"`
}

type Request struct {
	Name      string     `yaml:"name"`
	Method    string     `yaml:"method"`
	Path      string     `yaml:"path"`
	Responses []Response `yaml:"responses"`
}

type Response struct {
	Name     string   `yaml:"name"`
	Code     int      `yaml:"code"`
	Query    string   `yaml:"query"`
	Headers  []Header `yaml:"headers"`
	FilePath string   `yaml:"filePath"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func parseSetting(app_name string) *Setting {
	file_path := path.Join("data", app_name, "setting.yaml")
	data, err := os.ReadFile(file_path)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Parse the YAML data into the struct
	var setting Setting
	err = yaml.Unmarshal(data, &setting)
	if err != nil {
		log.Println(err)
		return nil
	}

	setting.Folder = app_name

	return &setting
}

func (s *Setting) loadResources(cacheManager *CacheManager) {
	// cache openapi files
	if s.SwaggerEnabled {
		for _, file := range []string{"openapi.json", "openapi.yml", "openapi.yaml"} {
			file_path := path.Join("data", s.Folder, file)
			data, err := os.ReadFile(file_path)
			if err == nil {
				cacheManager.update(file_path, data)
			}
		}
	}

	// cache response files of all requests
	for _, req := range s.Requests {
		for _, resp := range req.Responses {
			if resp.FilePath == "" {
				continue
			}
			data, err := os.ReadFile(resp.FilePath)
			if err != nil {
				log.Printf("Can not load resource req: %s resp: %s file %s \n", req.Name, resp.Name, resp.FilePath)
				panic(err)
			}
			cacheManager.update(resp.FilePath, data)
		}
	}
}
