package conf

import (
	"flag"
	"os"
)

func Example() {
	type Config struct {
		Port int    `default:"8080" env:"PORT" flag:",specify the port to listen on"`
		Host string `default:"localhost"`
		DB   struct {
			DSN string `default:"sqlite://:memory:"`
		}
	}

	var config Config
	// Sets the struct values using flags, environment variables and finally the default values defined on the struct
	if err := Load(&config, Auto()); err != nil {
		panic(err)
	}

	flag.Usage()
	// flag.Usage will output the following:
	// Usage of test:
	//   -db-dsn string
	//     	set db-dsn
	//   -host string
	//     	set host
	//   -port int
	//     	specify the port to listen on
}

func Example_auto() {
	type Config struct {
		Port int    `default:"8080" env:"PORT"`
		Host string `default:"localhost"`
		DB   struct {
			DSN string `default:"sqlite://:memory:"`
		}
	}
	var c Config

	// This will use any flags if they are provided, then environment variables and finally default values
	if err := Load(&c, Auto()); err != nil {
		panic(err)
	}

}

func Example_manual() {
	type Config struct {
		Port int    `default:"8080" env:"PORT"`
		Host string `default:"localhost"`
		DB   struct {
			DSN string `default:"sqlite://:memory:"`
		}
	}
	var c Config
	// This will use any flags if they are provided, then environment variables and finally default values
	if err := Load(&c, Flag(), Env(), Defaults()); err != nil {
		panic(err)
	}

}

// The default usage of the flag configurer
func ExampleFlag() {
	type Config struct {
		Port int    `flag:"port"`
		Host string `flag:"host"`
		DB   struct {
			DSN string `flag:"db-dsn"`
		}
	}
	// The following command line arguments will set the values on the struct:
	// -port 8080 -host localhost -db-dsn sqlite://:memory:
	var config Config
	if err := Load(&config, Flag()); err != nil {
		panic(err)
	}

}

// This is how we can customize how the flag configurer works.
func ExampleFlag_customization() {
	type Config struct {
		Port int    `flag:"port"`
		Host string `flag:"host"`
		DB   struct {
			DSN string `flag:"db-dsn"`
		}
	}
	// The following command line arguments will set the values on the struct:
	// -port 8080 -host localhost -db-dsn sqlite://:memory:
	var config Config
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	// This will set the flag names on the 'test' flag set
	if err := Load(&config, Flag(UseFlagSet(set))); err != nil {
		panic(err)
	}

	// This will not call flag.Parse()
	// if err := Load(&config, Flag(SkipParse())); err != nil {
	// 	panic(err)
	// }

}

func ExampleEnvFiles() {
	type Config struct {
		Port int    `env:"PORT"`
		Host string `env:"HOST"`
		DB   struct {
			DSN string `env:"DB_DSN"`
		}
	}

	var config Config
	// This will load the values from the .env file and ignores anything in os.Environ
	if err := Load(&config, EnvFiles(".env")); err != nil {
		panic(err)
	}

}

func ExampleEnv() {
	type Config struct {
		Port int    `env:"PORT"`
		Host string `env:"HOST"`
		DB   struct {
			DSN string `env:"DB_DSN"`
		}
	}

	// The following environment variables will set the values on the struct
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("DB_DSN", "sqlite://:memory:")
	var config Config
	if err := Load(&config, Env()); err != nil {
		panic(err)
	}

}
