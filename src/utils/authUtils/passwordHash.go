package authUtils

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareStrings(hashedPassword string, plainPassword []byte) (bool, error) {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	//if there's an error => (err == nil) == false
	// else  => (err == nil) == true
	return err == nil, err
}
