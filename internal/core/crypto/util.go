package crypto

func HashPassword(password string) (string, error) {
	hashFunc, err := NewArgon2ID()
	if err != nil {
		return "", err
	}
	hash, err := hashFunc.DeriveKey([]byte(password))
	if err != nil {
		return "", err
	}
	passwordHash := hash.String()
	return passwordHash, nil
}

func ComparePasswordHash(hash, password string) (bool, error) {
	kdf, existingHash, err := NewArgon2IDFromString(hash)
	if err != nil {
		return false, err
	}
	passwordHash, err := kdf.DeriveKey([]byte(password))
	if err != nil {
		return false, err
	}
	return kdf.Compare(existingHash, passwordHash.Bytes()), nil
}
