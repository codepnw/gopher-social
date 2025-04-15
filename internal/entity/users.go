package entity

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
	RoleID    int64  `json:"role_id"`
	Role      Role   `json:"role"`
}

type UserWithToken struct {
	*User
	Token string `json:"token"`
}

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) ComparePassword(text string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(text))
}
