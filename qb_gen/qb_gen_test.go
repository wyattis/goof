package qb_gen

import (
	"testing"

	"github.com/wyattis/goof/qb_gen/test_models"
)

func TestGenerate(t *testing.T) {
	err := Generate(Config{}, test_models.User{}, test_models.Comment{})
	if err != nil {
		t.Fatal(err)
	}
}
