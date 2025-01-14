package jwt

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ServiceName string `json:"service_name"`
	jwt.RegisteredClaims
}

func GenerateServiceToken() {
	// Секретный ключ для подписи токена
	secretKey := []byte("your_secret_key")

	// Создаем claims для токена
	claims := Claims{
		ServiceName: "service",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Ошибка при подписании токена: %v", err)
	}

	// Печатаем токен
	fmt.Printf("Сгенерированный JWT-токен: %s\n", signedToken)

	// Записываем токен в .env файл
	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("Ошибка при создании файла .env: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("JWT_TOKEN=%s\n", signedToken))
	if err != nil {
		log.Fatalf("Ошибка при записи в файл .env: %v", err)
	}

	fmt.Println("JWT-токен успешно сохранен в .env файл.")
}

// getJWTFromEnv читает JWT-токен из указанного .env файла
func GetJWTFromEnv(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	lines := string(data)
	for _, line := range strings.Split(lines, "\n") {
		if strings.HasPrefix(line, "JWT_TOKEN=") {
			return strings.TrimPrefix(line, "JWT_TOKEN="), nil
		}
	}

	return "", fmt.Errorf("JWT_TOKEN не найден в файле")
}
