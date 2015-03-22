package env_strings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ENV_STRINGS_KEY = "ENV_STRINGS"
	ENV_STRINGS_EXT = ".env"
)

type EnvStrings struct {
	envName string
	envExt  string
}

func NewEnvStrings(envName string, envExt string) *EnvStrings {
	if envName == "" {
		panic("env_strings: env name could not be nil")
	}

	return &EnvStrings{
		envName: envName,
		envExt:  envExt,
	}
}

func (p *EnvStrings) Execute(str string) (ret string, err error) {
	strConfigFiles := os.Getenv(p.envName)

	configFiles := strings.Split(strConfigFiles, ";")

	if strConfigFiles == "" || len(configFiles) == 0 {
		return str, nil
	}

	files := []string{}

	for _, confFile := range configFiles {
		var fi os.FileInfo
		if fi, err = os.Stat(confFile); err != nil {
			return
		}

		if fi.IsDir() {
			var dir *os.File
			if dir, err = os.Open(confFile); err != nil {
				return
			}

			var names []string
			if names, err = dir.Readdirnames(-1); err != nil {
				return
			}

			for _, name := range names {
				if ext := filepath.Ext(name); ext == p.envExt {
					filePath := strings.TrimRight(confFile, "/")
					files = append(files, filePath+"/"+name)
				}
			}
		} else {
			if ext := filepath.Ext(confFile); ext == p.envExt {
				files = append(files, confFile)
			}
		}
	}

	envs := map[string]map[string]interface{}{}

	for _, file := range files {
		var str []byte
		if str, err = ioutil.ReadFile(file); err != nil {

			return
		}

		env := map[string]interface{}{}
		if err = json.Unmarshal(str, &env); err != nil {
			return
		}

		envs[file] = env
	}

	allEnvs := map[string]interface{}{}

	for file, env := range envs {
		for envKey, envVal := range env {
			if _, exist := allEnvs[envKey]; exist {
				err = fmt.Errorf("env key of %s already exist, env file: %s", envKey, file)
				return
			} else {
				allEnvs[envKey] = envVal
			}
		}
	}

	var tpl *template.Template

	if tpl, err = template.New("env_strings").Parse(str); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, allEnvs); err != nil {
		return
	}

	ret = buf.String()

	if strings.Contains(ret, "<no value>") {
		err = fmt.Errorf("some env value did not exist, content: \n%s\n", ret)
		return
	}

	return
}

func Execute(str string) (ret string, err error) {
	envStrings := NewEnvStrings(ENV_STRINGS_KEY, ENV_STRINGS_EXT)
	return envStrings.Execute(str)
}
