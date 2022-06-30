package users

type UserRepository interface {
	FindAll() ([]User, error)
	FindOneByID(string) (*User, error)
	FindOneByEmail(string) (*User, error)
	Create(*User) (*User, error)
	UpdateOneByID(string, *User) error
	DeleteOneByID(string) error
}
