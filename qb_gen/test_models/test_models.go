package test_models

type User struct {
	Id   int
	Name string
}

type Comment struct {
	Id     int
	UserId int
	Body   string
}
