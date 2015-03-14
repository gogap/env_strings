ENV STRINGS
===========

Read Env from config and compile the values into string.

### About ENV JSON

Sometimes we need use some config string as following:

`db.conf`

```go
configStr:="127.0.0.1|123456|1000"
```

but, when we management more and more server and serivce, and if we need change the password or ip, it was a disaster.

So, we just want use config string like this.

`db.conf`

```json
configStr:="{{.host}}|{{.password}}|{{.timeout}}"
```

We use golang's template to replace values into the config while we execute the string.

first, we set the env config at `~/.bash_profile` or `~/.zshrc`, and the default key is `ENV_STRINGS` and the default file extention is `.env`, the value of `ENV_STRINGS ` could be a file or folder,it joined by`;`, it will automatic load all `*.env` files.

**Env**

```bash
export ENV_STRINGS ='~/playgo/test.env;~/playgo/test2.env'
```

or

```bash
export ENV_STRINGS ='~/playgo'
```


#### example program

```go
package main

import (
	"fmt"

	"github.com/gogap/env_strings"
)

func main() {
	if ret, err := env_strings.Execute("{{.host}}|{{.password}}|{{.timeout}}"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)
	}

	envStrings := env_strings.NewEnvStrings("ENV_KEY", ".env")

	if ret, err := envStrings.Execute("{{.host}}|{{.password}}|{{.timeout}}"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)
	}
}

```
 

`env1.json`

```json
{
	"host":"127.0.0.1",
	"password":"123456"
}
```


`env2.json`

```json
{
	"timeout":1000
}
```

**result:**

```bash
{127.0.0.1 123456 1000}
```
