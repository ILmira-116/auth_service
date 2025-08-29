package tests

import (
	"auth-service/tests/suite"
	"testing"
	"time"

	"github.com/ILmira-116/protos/gen/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	appID     = 1
	appSecret = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	// Генерация случайных данных
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	// Регистрация пользователя
	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err) // не смогли создать клиента дальше не продолжаем тест
	assert.NotEmpty(t, respReg.GetUserId())

	// Пользователь логинится
	respLogin, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	// получаем токен
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	// парсинг токена
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})

	// проходит ли токен валидацию
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	//сверяем поля токена
	uidClaim, ok := claims["user_id"].(float64)
	require.True(t, ok, "user_id claim must be float64")
	assert.Equal(t, respReg.GetUserId(), int64(uidClaim))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	// проверка времени истечения токена(с точностью до 1 сек)
	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}

// fail-кейс: пользователь регистрируется два раза
func TestRegisterLogin_DuplicateRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	// Первая регистрация — должна пройти
	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	// Вторая регистрация — ожидаем AlreadyExists
	respReg, err = st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	require.Nil(t, respReg) // безопасная проверка

	// Проверяем gRPC статус ошибки
	sts, ok := status.FromError(err)
	require.True(t, ok, "ошибка должна быть gRPC status")
	assert.Equal(t, codes.AlreadyExists, sts.Code())
	assert.Equal(t, "user already exists", sts.Message())
}
