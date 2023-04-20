package service_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	entity "github.com/ropel12/project-3/app/features/user"
	mocks "github.com/ropel12/project-3/app/features/user/mocks/repository"
	user "github.com/ropel12/project-3/app/features/user/service"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("user", func() {
	var Mock *mocks.UserRepo
	var UserService user.UserService
	var Depend dependcy.Depend
	var ctx context.Context
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		Mock = mocks.NewUserRepo(GinkgoT())
		UserService = user.NewUserService(Mock, Depend)

	})
	Context("User Login", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Eror", func() {
				err, _ := UserService.Login(ctx, entity.LoginReq{})
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Email Tidak terfadtar", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("Email not registered")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email not registered'", func() {
				_, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email not registered"))
			})
		})
		When("Password Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(&entity.User{Email: "satrio2@gmail.com", Password: "321"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'wrong password' ", func() {
				_, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Wrong password"))
			})
		})
		When("Berhasil Login", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", Password: "$2a$10$vu7o2Wl9LKyzTFkRDp7tc.VyoBB48nj97qyQjlgGCeQXJ067KZGQu"}
				data.ID = 1
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error", func() {
				uid, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).Should(BeNil())
				Expect(uid).To(Equal(1))
			})
		})

	})
	Context("User Register", func() {
		When("Request body kosong", func() {
			It("Akan Mengembalikan error", func() {
				err := UserService.Register(ctx, entity.RegisterReq{})
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Email sudah terdaftar", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com"}
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", Name: "satrio", Password: "123", Address: "bogor ct"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email already registered"))
			})
		})

		When("Passoword Terlalu panjang", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", Name: "satrio", Password: "eeewwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww", Address: "bogor ct"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Register failed"))

			})
		})
		When("Query database Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(errors.New("Failed to create account")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", Name: "satrio", Password: "eeeww", Address: "bogor ct"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to create account"))

			})
		})

		When("Query database Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", Name: "satrio", Password: "eeeww", Address: "bogor ct"})
				Expect(err).Should(BeNil())

			})
		})

	})

})
