package interfaces

type TokenProvider interface {
	Generate(id uint, username, role string) (string, error)
}
