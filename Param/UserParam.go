package Param

type UserParam struct {
	Username   string `form:"username"`
	Pwd        string `form:"pwd"`
	Email      string `form:"email"`
	VerifyCode string `form:"verifyCode"`
}
