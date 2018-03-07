package user

func (u *User) CreateUser() error {
	return UserTable.Create(u)
}
