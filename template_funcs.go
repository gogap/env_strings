package env_strings

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func newFuncMap() map[string]interface{} {
	m := make(map[string]interface{})
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
