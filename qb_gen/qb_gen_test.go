package qb_gen

import (
	"testing"

	"github.com/wyattis/goof/qb_gen/test_models"
)

func TestGenerateQb(t *testing.T) {
	err := Generate(Config{
		QueryBuilders: []any{test_models.User{}, test_models.Comment{}},
		Crud:          []any{test_models.User{}, test_models.Comment{}},
	})
	if err != nil {
		t.Fatal(err)
	}
}