package delivery

const (
	StatusUnauthorized = 401

	ResponseSuccessfulSignUp = "Successful sign up"
	ResponseSuccessfulSignIn = "Successful sign in"
	ResponseSuccessfulLogOut = "Successful log out"

	ErrUserNotExits     = "User with same email not exist"
	ErrUserAlreadyExist = "User with same email already exist"
	ErrWrongCredentials = "Uncorrected login or password"
	ErrUnauthorized     = "You unauthorized"
)
