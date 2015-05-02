package env_strings

type EnvStorage interface {
	FuncName() string
	Get(key string) (val string)
}
