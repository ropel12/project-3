package helper

import (
	"log"
	"time"

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

func GenerateEndTime(timee string, duration float32) string {
	t, err := time.Parse("2006-01-02 15:04:05", timee)
	if err != nil {
		return ""
	}
	minute := duration * 60

	return t.Add(time.Minute * time.Duration(int(minute))).Format("2006-01-02 15:04:05")
}
