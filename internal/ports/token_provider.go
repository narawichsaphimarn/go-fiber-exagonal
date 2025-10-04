package ports

type TokenProvider interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(token string) (string, error)
}
