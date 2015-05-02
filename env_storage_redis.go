package env_strings

import (
	"github.com/hoisie/redis"
)

type EnvStorageRedis struct {
	client redis.Client
	prefix string
}

func NewEnvStorageRedis(options map[string]interface{}) EnvStorage {
	storage := new(EnvStorageRedis)

	var addr string

	if v, exist := options["address"]; !exist {
		panic("option of address not exist")
	} else if strAddr, ok := v.(string); !ok {
		panic("option of address must be string")
	} else {
		addr = strAddr
	}

	var db float64
	if v, exist := options["db"]; !exist {
		db = 0
	} else if intDb, ok := v.(float64); !ok {
		panic("option of db must be int")
	} else {
		db = intDb
	}

	var password string
	if v, exist := options["password"]; !exist {
		password = ""
	} else if strPassword, ok := v.(string); !ok {
		panic("option of password must be string")
	} else {
		password = strPassword
	}

	var poolSize float64
	if v, exist := options["pool_size"]; !exist {
		poolSize = 0
	} else if intPoolSize, ok := v.(float64); !ok {
		panic("option of poolSize must be int")
	} else {
		poolSize = intPoolSize
	}

	var prefix string

	if v, exist := options["prefix"]; !exist {
		prefix = ""
	} else if strPrefix, ok := v.(string); !ok {
		panic("option of prefix must be string")
	} else {
		prefix = strPrefix
	}

	client := redis.Client{
		Addr:        addr,
		Db:          int(db),
		Password:    password,
		MaxPoolSize: int(poolSize),
	}

	storage.client = client
	storage.prefix = prefix

	return storage
}

func (p *EnvStorageRedis) FuncName() string {
	return "getv"
}

func (p *EnvStorageRedis) Get(key string) (val string) {
	if v, e := p.client.Get(p.prefix + "/" + key); e != nil {
		return ""
	} else {
		val = string(v)
	}
	return
}
