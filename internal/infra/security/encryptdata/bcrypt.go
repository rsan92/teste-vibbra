package encryptdata

import "golang.org/x/crypto/bcrypt"

type BCryptDataEncrypt struct {
}

func NewDataEncrypt() IEncrypt {
	return BCryptDataEncrypt{}
}

func (dataEncrypt BCryptDataEncrypt) GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (dataEncrypt BCryptDataEncrypt) VerifyHash(passwordString, hashString string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashString), []byte(passwordString))
}
