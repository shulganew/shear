package config

// send values through middleware in context
type CtxConfig struct {
	userID    string
	isNewUser bool
}

func NewCtxConfig(userID string, isNewUser bool) CtxConfig {

	return CtxConfig{userID: userID, isNewUser: isNewUser}
}

func (c CtxConfig) GetUserID() string {
	return c.userID
}

func (c CtxConfig) IsNewUser() bool {
	return c.isNewUser
}