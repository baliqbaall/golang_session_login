package entities

type User struct {
	Id              int64
	NamaLengkap     string `validate:"required" label:"Fullname"`
	Email           string `validate:"required,email,isunique=users-email"`
	Username        string `validate:"required,gte=3,isunique=users-username"`
	Password        string `validate:"required,gte=6"`
	ConfirmPassword string `validate:"required,eqfield=Password" label:"Confirm Password"`
}
