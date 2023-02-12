package encryptdata

type IEncrypt interface {
	GenerateHash(string) ([]byte, error)
	VerifyHash(string, string) error
}
