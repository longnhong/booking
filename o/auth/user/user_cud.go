package user

func (u *User) CreateUser() error {
	var psd, err = Password(u.Password).GererateHashedPassword()
	if err != nil {
		return err
	}
	u.Password = psd
	if u.Role == 0 {
		u.Role = ROLE_USER
	}
	return UserTable.Create(u)
}
