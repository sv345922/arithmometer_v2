package loginUser

import (
	"context"
	"errors"

	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/jwToken"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
)

func getAccess(ctx context.Context, ws *wSpace.WorkingSpace, data entities.UserData) (string, error) {
	// Найти id пользователя в БД
	pass, err := findUserPass(ctx, ws.DB, data.Name)
	if err != nil {
		return "", err
	}
	// сопоставить пароль
	if pass != data.Password {
		return "", errors.New("invalid password")
	}
	// Создать jwt и вернуть его
	token, err := jwToken.CreateJWT(data.Name, 5)
	if err != nil {
		return "", err
	}
	return token, err
}
