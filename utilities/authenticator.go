package utilities

type Authenticator struct {
	Passwords map[string]string
}

func NewAuthenticator(defaultUser, defaultPassword string) *Authenticator {
	var auth Authenticator
	auth.Passwords = make(map[string]string)
	auth.Passwords[defaultUser] = defaultPassword
	return &auth
}

func (auth Authenticator) Authenticate(username, password string) bool {
	rv := false

	if hash, ok := auth.Passwords[username]; ok {
		if password == hash {
			rv = true
		}
	}

	return rv
}

// END OF SOURCE
