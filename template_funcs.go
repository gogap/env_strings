package env_strings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"
)

type TemplateFuncs struct {
	funcMap template.FuncMap
}

func NewTemplateFuncs() *TemplateFuncs {
	return &TemplateFuncs{funcMap: basicFuncs()}
}

func (p *TemplateFuncs) All() template.FuncMap {
	return p.funcMap
}

func (p *TemplateFuncs) Register(name string, function interface{}) (err error) {
	if function == nil {
		err = errors.New("function could not be nil")
		return
	} else if funcType := reflect.TypeOf(function); funcType.Kind() != reflect.Func {
		err = errors.New("function is not a Func kind")
		return
	} else if name == "" {
		err = errors.New("name could not be empty")
		return
	}

	if _, exist := p.funcMap[name]; exist {
		panic("func name of " + name + " already exist")
	}

	p.funcMap[name] = function

	return
}

func basicFuncs() template.FuncMap {
	m := make(template.FuncMap)
	m["base"] = path.Base
	m["abs"] = filepath.Abs
	m["dir"] = filepath.Dir
	m["pwd"] = os.Getwd
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["getenv"] = os.Getenv
	m["join"] = strings.Join
	m["localtime"] = time.Now
	m["utc"] = time.Now().UTC
	m["pid"] = os.Getpid
	m["httpGet"] = httpGet
	m["envIfElse"] = envIfElse

	return m
}

func UnmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func UnmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func httpGet(url string) (ret string, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	ret = string(body)

	return
}

func envIfElse(envName, equalValue, trueValue, falseValue string) string {
	if os.Getenv(envName) == equalValue {
		return trueValue
	} else {
		return falseValue
	}
}
