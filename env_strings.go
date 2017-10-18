package env_strings

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ENV_STRINGS_KEY = "ENV_STRINGS"
	ENV_STRINGS_EXT = ".env"

	ENV_STRINGS_CONF = "/etc/env_strings.conf"

	ENV_STRINGS_CONFIG_KEY = "ENV_STRINGS_CONF"
)

const (
	STORAGE_REDIS = "redis"
)

type EnvStringConfig struct {
	Storages []StorageConfig `json:"storages"`
}

type StorageConfig struct {
	Engine  string                 `json:"engine"`
	Options map[string]interface{} `json:"options"`
}

type option func(envStrings *EnvStrings)

type EnvStrings struct {
	envName   string
	envExt    string
	tmplFuncs *TemplateFuncs

	configFile string

	envConfig EnvStringConfig
}

func FuncMap(name string, function interface{}) option {
	return func(e *EnvStrings) {
		e.RegisterFunc(name, function)
	}
}

func EnvStringsConfig(fileName string) option {
	return func(e *EnvStrings) {
		e.configFile = fileName
	}
}

func NewEnvStrings(envName string, envExt string, opts ...option) *EnvStrings {
	if envName == "" {
		panic("env_strings: env name could not be empty")
	}

	envStrings := &EnvStrings{
		envName:    envName,
		envExt:     envExt,
		configFile: ENV_STRINGS_CONF,
		tmplFuncs:  NewTemplateFuncs(),
	}

	if opts != nil && len(opts) > 0 {
		for _, opt := range opts {
			opt(envStrings)
		}
	}

	envStringsConf := os.Getenv(ENV_STRINGS_CONFIG_KEY)
	if envStringsConf != "" {
		envStrings.configFile = envStringsConf
	}

	if envStrings.configFile != "" {
		if err := envStrings.loadConfig(envStrings.configFile); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			} else {
				return envStrings
			}
		}

		if envStrings.envConfig.Storages != nil {
			for _, storageConf := range envStrings.envConfig.Storages {
				switch storageConf.Engine {
				case STORAGE_REDIS:
					{
						extFucnRedis := NewExtFuncsRedis(storageConf.Options)
						redisFuncs := extFucnRedis.GetFuncs()

						if redisFuncs == nil {
							panic("ext funcs of redis is nil")
						}

						for funcName, fn := range redisFuncs {
							if err := envStrings.RegisterFunc(funcName, fn); err != nil {
								panic(err)
							}
						}
					}
				default:
					{
						panic("unknown storage type")
					}
				}
			}
		}
	}
	return envStrings
}

func (p *EnvStrings) Execute(str string) (ret string, err error) {
	return p.ExecuteWith(str, nil)
}

func (p *EnvStrings) ExecuteWith(str string, envValues map[string]interface{}) (ret string, err error) {
	strEnvFiles := os.Getenv(p.envName)

	envFiles := strings.Split(strEnvFiles, ";")

	if len(envFiles) == 1 && len(envFiles[0]) == 0 {
		envFiles = nil
	}

	if envValues == nil {
		envValues = make(map[string]interface{})
	}

	for _, envFile := range envFiles {

		prefix := filepath.Base(envFile)

		prefix = filepath.Dir(prefix)

		err = p.loadEnv(prefix, []string{envFile}, envValues)

		if err != nil {
			return
		}
	}

	var tpl *template.Template

	if tpl, err = template.New("tmpl:" + p.envName).Funcs(p.tmplFuncs.GetFuncMaps(p.envName)).Option("missingkey=error").Parse(str); err != nil {
		return
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, envValues); err != nil {
		return
	}

	ret = buf.String()

	return
}

func Execute(str string) (ret string, err error) {
	envStrings := NewEnvStrings(ENV_STRINGS_KEY, ENV_STRINGS_EXT)
	return envStrings.Execute(str)
}

func ExecuteWith(str string, envValues map[string]interface{}) (ret string, err error) {
	envStrings := NewEnvStrings(ENV_STRINGS_KEY, ENV_STRINGS_EXT)
	return envStrings.ExecuteWith(str, envValues)
}

func (p *EnvStrings) RegisterFunc(name string, function interface{}) (err error) {
	return p.tmplFuncs.Register(name, function)
}

func (p *EnvStrings) FuncUsageStatic() map[string][]FuncStaticItem {
	return funcStatics
}

func (p *EnvStrings) loadEnv(prefix string, files []string, envs map[string]interface{}) (err error) {

	for _, path := range files {

		var fi os.FileInfo
		fi, err = os.Stat(path)
		if err != nil {
			return
		}

		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		if fi.IsDir() {

			baseName := filepath.Base(path)

			if strings.HasPrefix(baseName, ".") {
				continue
			}

			var fis []os.FileInfo
			fis, err = ioutil.ReadDir(path)

			if err != nil {
				return
			}

			var nextfiles []string

			for _, f := range fis {
				nextfiles = append(nextfiles, filepath.Join(path, f.Name()))
			}

			var nextENVs map[string]interface{}

			if envs == nil {
				envs = make(map[string]interface{})
			}

			preEnvs, exist := envs[baseName]
			if !exist {
				nextENVs = make(map[string]interface{})
				envs[baseName] = nextENVs
			} else {
				nextENVs = preEnvs.(map[string]interface{})
			}

			err = p.loadEnv(prefix, nextfiles, nextENVs)
			if err != nil {
				return
			}

			continue
		}

		baseName := strings.TrimSuffix(fi.Name(), p.envExt)

		fileEnvs, err := p.loadEnvFile(path)

		if err != nil {
			return err
		}

		if envs == nil {
			envs = make(map[string]interface{})
		}

		envs[baseName] = fileEnvs
	}

	return
}

func (p *EnvStrings) loadEnvFile(filename string) (ret map[string]interface{}, err error) {

	if ext := filepath.Ext(filename); ext == p.envExt {

		var data []byte
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}

		r := make(map[string]interface{})

		err = json.Unmarshal(data, &r)
		if err != nil {
			return
		}

		ret = r
	}

	return
}

func (p *EnvStrings) loadConfig(fileName string) (err error) {
	if _, err = os.Stat(fileName); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	conf := EnvStringConfig{}

	if err = json.Unmarshal(data, &conf); err != nil {
		return
	}

	p.envConfig = conf

	return
}
