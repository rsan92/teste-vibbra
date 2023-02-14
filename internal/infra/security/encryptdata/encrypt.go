package encryptdata

type IEncrypt interface {
	GenerateHash(string) ([]byte, error)
	VerifyHash(passwordString, hashString string) error
}
