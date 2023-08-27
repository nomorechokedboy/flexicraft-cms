package hasher

type Hasher interface {
	Hash(string) (*string, error)
	Verify(string, string) (*bool, error)
}
