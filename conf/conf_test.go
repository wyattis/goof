package conf

// func TestEnv(t *testing.T) {
// 	os.Setenv("GOOS", runtime.GOOS)
// 	os.Setenv("GOARCH", runtime.GOARCH)
// 	os.Setenv("GOOT", "wow")
// 	type config struct {
// 		GOOS   string
// 		GOARCH string
// 		Goot   string
// 	}
// 	conf := config{}
// 	c := New(Env())
// 	if err := c.Apply(&conf); err != nil {
// 		t.Error(err)
// 	}
// 	expected := config{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, Goot: "wow"}
// 	if !reflect.DeepEqual(conf, expected) {
// 		t.Errorf("Expected %+v, but got %+v", expected, conf)
// 	}
// }

// func TestDotEnv(t *testing.T) {
// 	os.Clearenv()
// 	data := []byte(fmt.Sprintf("GOOS=%s\nGOARCH=%s\nNESTED_HELLO=world\n", runtime.GOOS, runtime.GOARCH))
// 	if err := os.WriteFile(".env.test", data, os.ModePerm); err != nil {
// 		t.Error(err)
// 	}
// 	defer os.Remove(".env.test")
// 	type config struct {
// 		GOOS   string
// 		GOARCH string
// 		Nested struct {
// 			Hello string
// 		}
// 	}
// 	conf := config{}
// 	c := New(EnvFiles(".env.test"))
// 	if err := c.Apply(&conf); err != nil {
// 		t.Error(err)
// 	}
// 	expected := config{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, Nested: struct{ Hello string }{"world"}}
// 	if conf.GOARCH != expected.GOARCH || conf.GOOS != expected.GOOS || conf.Nested.Hello != expected.Nested.Hello {
// 		t.Errorf("Expected %+v, but got %+v", expected, conf)
// 	}
// }

// func TestFlag(t *testing.T) {
// 	type config struct {
// 		Addr   string
// 		Custom string `flag:"another"`
// 	}
// 	var cases = []struct {
// 		args []string
// 		res  config
// 	}{
// 		{[]string{"-addr", ":80"}, config{Addr: ":80"}},
// 		{[]string{"-another", "world"}, config{Custom: "world"}},
// 	}
// 	for _, c := range cases {
// 		val := config{}
// 		conf := New(Flag(c.args))
// 		if err := conf.Apply(&val); err != nil {
// 			t.Error(err)
// 		}
// 		if !reflect.DeepEqual(val, c.res) {
// 			t.Errorf("Expected %+v, but got %+v", c.res, val)
// 		}
// 	}
// }

// func TestFlagDefaults(t *testing.T) {
// 	type config struct {
// 		Addr   string `default:":80"`
// 		Custom string `flag:"another"`
// 		Nested struct {
// 			Hello string
// 		}
// 	}
// 	var cases = []struct {
// 		args []string
// 		res  config
// 	}{
// 		{[]string{"-addr", ":8080"}, config{Addr: ":8080"}},
// 		{[]string{"-another", "world"}, config{Addr: ":80", Custom: "world"}},
// 		{[]string{"-nested-hello", "world"}, config{Addr: ":80", Nested: struct{ Hello string }{"world"}}},
// 	}
// 	for _, c := range cases {
// 		val := config{}
// 		conf := New(Flag(c.args), Defaults())
// 		if err := conf.Apply(&val); err != nil {
// 			t.Error(err)
// 		}
// 		if !reflect.DeepEqual(val, c.res) {
// 			t.Errorf("Expected %+v, but got %+v", c.res, val)
// 		}
// 	}
// }

// func TestFailure(t *testing.T) {
// 	type config struct {
// 		Int int
// 	}
// 	val := config{}
// 	c := New(Flag([]string{"-int asdf"}))
// 	if err := c.Apply(&val); err == nil {
// 		t.Error("expected error parsing integer")
// 	}
// }
