package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	repositorymock "GinStart/Repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	passwd := []byte("123123aaa")
	encrypt, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypt))
	err = bcrypt.CompareHashAndPassword(encrypt, []byte("123123aaa"))
	assert.NoError(t, err)
}

func Test_userService_Login(t *testing.T) {
	TestCases := []struct {
		name         string
		mock         func(ctrl *gomock.Controller) Repository.UserRepository
		ctx          context.Context
		email        string
		passwd       string
		ExpectedUser Domain.User
		ExpectedErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) Repository.UserRepository {
				repo := repositorymock.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(Domain.User{
						Id:       1,
						Password: "$2a$10$CWU7EDtqkv8vgT21nreNYu1CqFt4tsssbDjxF7sViDUfumLh6A0nq",
					}, nil)

				return repo
			},
			ctx:    context.Background(),
			email:  "123@qq.com",
			passwd: "123123aaa",
			ExpectedUser: Domain.User{
				Id:       1,
				Password: "$2a$10$CWU7EDtqkv8vgT21nreNYu1CqFt4tsssbDjxF7sViDUfumLh6A0nq",
			},
			ExpectedErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) Repository.UserRepository {
				repo := repositorymock.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(Domain.User{}, Repository.ErrUserNotFound)

				return repo
			},
			ctx:          context.Background(),
			email:        "123@qq.com",
			passwd:       "123123aaa",
			ExpectedUser: Domain.User{},
			ExpectedErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) Repository.UserRepository {
				repo := repositorymock.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(Domain.User{}, errors.New("系统错误"))

				return repo
			},
			ctx:          context.Background(),
			email:        "123@qq.com",
			passwd:       "123123aaa",
			ExpectedUser: Domain.User{},
			ExpectedErr:  errors.New("系统错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) Repository.UserRepository {
				repo := repositorymock.NewMockUserRepository(ctrl)

				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(Domain.User{
						Id:       1,
						Password: "$2a$10$CWU7EDtqkv8vgT21nreNYu1CqFt4tsssbDjxF7sViDUfumLh6A0nq",
					}, nil)

				return repo
			},
			ctx:          context.Background(),
			email:        "123@qq.com",
			passwd:       "123123",
			ExpectedUser: Domain.User{},
			ExpectedErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.Login(tc.ctx, tc.email, tc.passwd)
			assert.Equal(t, tc.ExpectedErr, err)
			assert.Equal(t, tc.ExpectedUser, user)
		})
	}
}
