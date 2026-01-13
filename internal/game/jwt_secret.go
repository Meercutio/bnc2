package game

var jwtSecretValue []byte

func SetJWTSecret(secret []byte) {
	jwtSecretValue = secret
}

func jwtSecret() []byte {
	return jwtSecretValue
}
