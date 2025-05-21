# env - Go Environment Variable Manager

`env` is a simple Go package that provides an environment variable manager. It supports loading environment variables from a `.env` file and the system environment, and provides convenient methods to retrieve values in different data types.

## Installation
To install the `env` package, use the following command:
```sh
go get github.com/azhai/go-env
```

## Usage

1. Create a .env File
   First, create a .env file in your project root directory with the following content:

```
# DATABASE_TYPE = "postgres"
DATABASE_URL = "postgres://${DEV_USER}:${DEV_PASSWORD}@127.0.0.1/test?sslmode=disable"

DB_HOST = localhost
DB_PORT = 3306
DB_USER = admin
DB_PASSWORD = "secret"
ENABLE_FEATURE = true
```

2. Import and Use the env Package
   Here is a simple example demonstrating how to use the env package:

```go
package main

import (
    "fmt"
    "github.com/azhai/go-env"
)

func main() {
    // Create a new Env instance and load variables from the default .env file
    e := env.New()

    // Get a string value
    dbHost := e.GetStr("DB_HOST", "127.0.0.1")
    fmt.Printf("DB Host: %s\n", dbHost)

    // Get an integer value
    dbPort := e.GetInt("DB_PORT", 3306)
    fmt.Printf("DB Port: %d\n", dbPort)

    // Get a boolean value
    enableFeature := e.GetBool("ENABLE_FEATURE", false)
    fmt.Printf("Enable Feature: %v\n", enableFeature)

    // Get a non-existent variable with a fallback value
    nonExistent := e.GetStr("NON_EXISTENT", "fallback_value")
    fmt.Printf("Non-existent Variable: %s\n", nonExistent)
}
```

3. Load from a Custom File
   If you want to load environment variables from a custom file, you can use the NewWithFile function:

```go
package main

import (
    "fmt"
    "os"
	"strings"
    "github.com/azhai/go-env"
)

// isRunTest is testing mode or not
func isRunTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func main() {
    // Create a new Env instance and load variables from a custom file
    var e *env.Env
    if isRunTest() {
		e = env.NewWithFile(".env.testing")
	} else {
		e = env.New()
	}

    // Do this in your shell: export DEV_USER=test DEV_PASSWORD=test
    // Show the DSN, replace the variables like DB_USER/DB_PASSWORD in it
    fmt.Printf("DSN: %s\n", e.Get("DATABASE_URL"))

	// Get an environment variable
	// Additionally, when e.Lookup() is called, it will attempt to load the environment variable from
	// the system environment variables, if it was not in the `.env` file.
	var port int
	if dbPort, ok := e.Lookup("DB_PORT"); ok {
		port = dbPort.Int()
	}
	// or
	// port := e.GetInt("DB_PORT")

	if port > 0 && port <= 65535 {
		fmt.Println("PORT:", port)
	} else {
		fmt.Println("PORT not set")
    }

}
```

## API Documentation

**Env Type**

The Env type represents an environment variable manager. It stores environment variables in a map and supports loading from a file and the system environment.

**New() \*Env**

Create a new Env instance and loads environment variables from the default .env file.

**NewWithFile(filename string) \*Env**

Create a new Env instance and loads environment variables from the specified file. If the file cannot be opened or read, the error is ignored.

**(e \*Env) Lookup(key string) (Entry, bool)**

Searche for an environment variable by key. It first checks the internal storage. If not found, it checks the system environment variables. If the variable is found, it returns the Entry and true; otherwise, it returns an empty Entry and false. If the variable is found in the system environment, it is added to the internal storage.

**(e \*Env) Load(reader io.ReadCloser, err error) error**

Read environment variables from a reader and stores them in the internal storage.

**(e \*Env) Get(key string) string**

Retrieve and expand the value of an environment variable by key.

**(e \*Env) GetStr(key string, fallback ...string) string**

Retrieve the string value of an environment variable by key. If the variable is not found, it returns the fallback value.

**(e \*Env) GetInt(key string, fallback ...int) int**

Retrieve the integer value of an environment variable by key. If the variable is not found or the value cannot be converted to an integer, it returns the fallback value.

**(e \*Env) GetInt64(key string, fallback ...int64) int64**

Retrieve the 64-bit integer value of an environment variable by key. If the variable is not found or the value cannot be converted to a 64-bit integer, it returns the fallback value.

**(e \*Env) GetBool(key string, fallback ...bool) bool**

Retrieve the boolean value of an environment variable by key. It supports "yes", "no", "true", and "false" as valid boolean values. If the variable is not found or the value cannot be converted to a boolean, it returns the fallback value.

**Entry Type**

The Entry type represents an environment variable entry with a key and a value.

**(e \*Entry) Str() string**

Return the string value of the entry.

**(e \*Entry) Int() int**

Convert the entry's value to an integer. If the conversion fails, it returns 0.

**(e \*Entry) Int64() int64**

Convert the entry's value to a 64-bit integer. If the conversion fails, it returns 0.

**(e \*Entry) Bool() bool**

Convert the entry's value to a boolean. It supports "yes", "no", "true", and "false" as valid boolean values. If the conversion fails, it returns false.

## License

Authors:
- [Nathan B. Crocker](https://github.com/nathanbcrocker)
- [Ryan Liu "azhai"](https://github.com/azhai)

Blog:
: [Simple Configuration Management in Go](https://medium.com/checker-engineering/simple-configuration-management-in-go-c4fe5db2f82e
)

This project is licensed under the GPLv3 License.
