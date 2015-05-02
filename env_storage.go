package env_strings

type EnvStorage interface {
	FuncName() string
	Get(key string, defaultVal ...string) (val string)
}
