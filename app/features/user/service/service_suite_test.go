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
		When("Request Body Tidak ada", func() {
			It("Akan Mengembalikan Eror", func() {
				err, _ := UserService.Login(ctx, entity.LoginReq{})
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Email Tidak terfadtar", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("Email tidak terdaftar")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email tidak terdaftar'", func() {
				_, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email tidak terdaftar"))
			})
		})
		When("Password Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(&entity.User{Email: "satrio2@gmail.com", Password: "321"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'crypto/bcrypt: hashedSecret too short to be a bcrypted password' ", func() {
				_, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Password Salah"))
			})
		})
		When("Jika Berhasil Login", func() {
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

})
