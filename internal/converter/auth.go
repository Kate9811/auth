// internal/converter/converter.go
package converter

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	//"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Denis/project_auth/internal/model"
	desc "github.com/Denis/project_auth/pkg/user_v1"
)

// Объявите свою ошибку
var (
	ErrPasswordsNotMatch = errors.New("passwords do not match")
)

// ToAuthInfoFromCreate конвертирует CreateRequest в AuthInfo
func ToAuthInfoFromCreate(req *desc.CreateRequest) (*model.AuthInfo, error) {
	// Валидация: пароли должны совпадать
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, ErrPasswordsNotMatch
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Конвертируем enum Role в string
	roleStr := ToModelRole(req.GetRole())

	return &model.AuthInfo{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: string(hash),
		Role:         roleStr,
	}, nil
}

// ToAuthInfoFromUpdate конвертирует UpdateRequest в AuthInfo
func ToAuthInfoFromUpdate(req *desc.UpdateRequest) *model.AuthInfo {
	info := &model.AuthInfo{}

	// Для google.protobuf.StringValue нужно использовать GetValue()
	if name := req.GetName(); name != nil {
		info.Name = name.GetValue()
	}

	if email := req.GetEmail(); email != nil {
		info.Email = email.GetValue()
	}

	return info
}

// AuthToGetResponse конвертирует Auth в GetResponse
func AuthToGetResponse(auth *model.Auth) *desc.GetResponse {
	// Конвертируем string role в enum
	role := ToDescRole(auth.Info.Role)

	resp := &desc.GetResponse{
		Id:        auth.ID,
		Name:      auth.Info.Name,
		Email:     auth.Info.Email,
		Role:      role,
		CreatedAt: timestamppb.New(auth.CreatedAt),
	}

	// UpdatedAt может быть NULL
	if auth.UpdatedAt.Valid {
		resp.UpdatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	return resp
}

// AuthToCreateResponse конвертирует ID в CreateResponse
func AuthToCreateResponse(id int64) *desc.CreateResponse {
	return &desc.CreateResponse{Id: id}
}

// // ToModelRole конвертирует protobuf Role в string
// func ToModelRole(role desc.Role) string {
// 	switch role {
// 	case desc.Role_ADMIN:
// 		return "ADMIN"
// 	default:
// 		return "USER"
// 	}
// }

// // ToDescRole конвертирует string в protobuf Role
// func ToDescRole(roleStr string) desc.Role {
// 	switch roleStr {
// 	case "ADMIN":
// 		return desc.Role_ADMIN
// 	default:
// 		return desc.Role_USER
// 	}
// }

// ToModelRole конвертирует protobuf Role в string
func ToModelRole(role desc.Role) string {
	switch role {
	case desc.Role_ADMIN:
		return "admin" // ← возможно нужно маленькими буквами
	default:
		return "user" // ← возможно нужно маленькими буквами
	}
}

// ToDescRole конвертирует string в protobuf Role
func ToDescRole(roleStr string) desc.Role {
	// Приводим к нижнему регистру для сравнения
	switch strings.ToLower(roleStr) {
	case "admin":
		return desc.Role_ADMIN
	default:
		return desc.Role_USER
	}
}
