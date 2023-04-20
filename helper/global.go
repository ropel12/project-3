package helper

import (
	"log"

	"github.com/golang-jwt/jwt"

	depedency "github.com/ropel12/project-3/config/dependcy"
	"golang.org/x/crypto/bcrypt"
)

func GetUid(token *jwt.Token) int {
	parse := token.Claims.(jwt.MapClaims)
	id := int(parse["id"].(float64))

	return id
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(passhash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password))
}

func GenerateJWT(id int, dp depedency.Depend) string {
	var informasi = jwt.MapClaims{}
	informasi["id"] = id
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, informasi)
	resultToken, err := rawToken.SignedString([]byte(dp.Config.JwtSecret))
	if err != nil {
		log.Println("generate jwt error ", err.Error())
		return ""
	}
	return resultToken
}
