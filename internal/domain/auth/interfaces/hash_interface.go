package interfaces

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type TokenProvider interface {
	Generate(id uint, username, role string) (string, error)
}
