package ports

type Crypt interface {
	GenerateFromPassword([]byte, int) ([]byte, error)
	CompareHashAndPassword([]byte, []byte) error
}
