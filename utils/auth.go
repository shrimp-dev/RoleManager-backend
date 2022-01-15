package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"drinkBack/models"
	"encoding/base64"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateRandomSalt(saltSize int) ([]byte, error) {
	// Generate salt
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt[:])

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func HashPassword(password string, salt []byte) (string, error) {
	sha512Hasher := sha512.New()

	if _, err := sha512Hasher.Write([]byte(password)); err != nil {
		return "", err
	}

	hashedPasswordBytes := sha512Hasher.Sum(salt)
	base64EncodedPasswordHash := base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash, nil
}

func MatchPassword(hashedPassword string, currPassword string, salt []byte) (bool, error) {
	currPasswordHash, err := HashPassword(currPassword, salt)
	if err != nil {
		return false, err
	}

	return hashedPassword == currPasswordHash, nil
}

func ValidatePassword(password []byte) (bool, error) {
	//regex -> (?!.*?[=?<>()'"\/\&]).{8,20}
	if len(password) < 8 && len(password) > 20 {
		return false, errors.New("invalid password size")
	}
	// Contains special char
	schar, err := regexp.Compile(`[!|@|#|\$|%|\*|-|_|\+|=]`)
	if err != nil {
		return false, err
	}
	// Upper case char
	uchar, err := regexp.Compile(`[A-Z]+`)
	if err != nil {
		return false, err
	}

	lchar, err := regexp.Compile(`[a-z]+`)
	if err != nil {
		return false, err
	}

	nchar, err := regexp.Compile(`\d+`)
	if err != nil {
		return false, err
	}

	// Contains script text or space
	soschar, err := regexp.Compile(`[!|@|#|\$|%|\*|-|_|\+|=]`)
	if err != nil {
		return false, err
	}
	// Check regex
	if schar.Match(password) &&
		uchar.Match(password) &&
		lchar.Match(password) &&
		nchar.Match(password) &&
		soschar.Match(password) {
		return true, nil
	}
	return false, nil
}

// jwt
const (
	AUTH   string = "AUTH_SECRET"
	INVITE string = "INVITE_SECRET"
)

func GenerateAuthenticationToken(id string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv(secret)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAuthenticationToken(tokenString string, secret string, dataCollector *models.AccessTokenClaims) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return false, nil //errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(os.Getenv(secret)), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("this token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	dataCollector.Id = claims["_id"].(string)

	return true, nil
}
