package conf

import (
	"os"
	"testing"
)

type envTestStrings struct {
	ABCD   string
	Secret struct {
		Test struct {
			AAA string
		}
	}
	Tagged int `env:"AAA_NAME"`
}

func TestZenvEnviron(t *testing.T) {
	os.Setenv("ABCD", "abcd")
	os.Setenv("SECRET_TEST_AAA", "aaa")
	os.Setenv("AAA_NAME", "123")
	s := envTestStrings{}
	c := envConfigurer{}
	if err := c.Apply(&s); err != nil {
		t.Error(err)
	}
	if s.ABCD != "abcd" {
		t.Errorf("Expected %s, but got %s", "abcd", s.ABCD)
	}
	if s.Secret.Test.AAA != "aaa" {
		t.Errorf("Expected %s, but got %s", "aaa", s.Secret.Test.AAA)
	}
	if s.Tagged != 123 {
		t.Errorf("Expected %d, but got %d", 123, s.Tagged)
	}
}

type envTestFileStruct struct {
	Hello string
}

func TestZenvExistingFiles(t *testing.T) {
	s := envTestFileStruct{}
	c := envConfigurer{filepaths: []string{"testdata/test.env"}}
	if err := c.Apply(&s); err != nil {
		t.Error(err)
	}
	if s.Hello != "WORLD" {
		t.Errorf("Expected %s, but got %s", "WORLD", s.Hello)
	}
	s = envTestFileStruct{}
	c = envConfigurer{filepaths: []string{"testdata/test.env.dev"}}
	if err := c.Apply(&s); err != nil {
		t.Error(err)
	}
	if s.Hello != "NOONE" {
		t.Errorf("Expected %s, but got %s", "NOONE", s.Hello)
	}
	s = envTestFileStruct{}
	c = envConfigurer{filepaths: []string{"testdata/test.env", "testdata/test.env.dev"}}
	if err := c.Apply(&s); err != nil {
		t.Error(err)
	}
	if s.Hello != "NOONE" {
		t.Errorf("Expected %s, but got %s", "NOONE", s.Hello)
	}
}

func TestZenvNonexistingFiles(t *testing.T) {
	// By default it ignores non-existing files
	s := envTestFileStruct{}
	c := envConfigurer{filepaths: []string{"testdata/test.env", "testdata/test.env.dev", "testdata/test.env.nonexisting", "fakedir/one.env"}}
	if err := c.Apply(&s); err != nil {
		t.Error(err)
	}
	if s.Hello != "NOONE" {
		t.Errorf("Expected %s, but got %s", "NOONE", s.Hello)
	}
	c = envConfigurer{filepaths: []string{"testdata/test.env.nonexisting"}, fileMustExist: true}
	if err := c.Apply(&s); err == nil {
		t.Error("Expected error, but got nil")
	}
	c = envConfigurer{filepaths: []string{"fakedir/one.env"}, fileMustExist: true}
	if err := c.Apply(&s); err == nil {
		t.Error("Expected error, but got nil")
	}
}
