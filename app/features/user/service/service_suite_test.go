package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	entity2 "github.com/ropel12/project-3/app/entities"
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
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(&entity2.User{Email: "satrio2@gmail.com", Password: "321"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'wrong password' ", func() {
				_, err := UserService.Login(ctx, entity.LoginReq{Email: "satrio2@gmail.com", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Wrong password"))
			})
		})
		When("Berhasil Login", func() {
			BeforeEach(func() {
				data := &entity2.User{Email: "satrio2@gmail.com", Password: "$2a$10$vu7o2Wl9LKyzTFkRDp7tc.VyoBB48nj97qyQjlgGCeQXJ067KZGQu"}
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
				data := &entity2.User{Email: "satrio2@gmail.com"}
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", Name: "satrio", Password: "123", Address: "bogor ct"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email already registered"))
			})
		})

		When("Password Terlalu panjang (melebihi 72 char)", func() {
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

	Context("User Update", func() {
		When("Password baru terlalu panjang (melebihi 72 char)", func() {
			It("Akan Mengembalikan error", func() {
				var file multipart.File
				_, err := UserService.Update(ctx, entity.UpdateReq{Password: "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww"}, file)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Email baru sudah terdaftar pada database", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity2.User{Email: "satrio@gmail.com"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email already registered'", func() {
				var file multipart.File
				_, err := UserService.Update(ctx, entity.UpdateReq{Password: "wwwwwwww", Email: "satrio@gmail.com"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email already registered"))
			})
		})
		When("Terdapat request gambar yang tidak sesuai dengan format gambar", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Failed to upload image' ", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				_, err := UserService.Update(ctx, entity.UpdateReq{Password: "satrio2323223", Email: "satrio3@gmail.com", Image: "gambar.php"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to upload image"))
			})
		})

		When("Kesalahan Database", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("Error db")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				_, err := UserService.Update(ctx, entity.UpdateReq{Password: "satrio2323223", Email: "satrio44@gmail.com", Image: "gambar.jpg"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Error db"))
			})
		})
		When("Berhasil mengupdate profile", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(&entity2.User{Email: "satrio44@gmail.com"}, nil).Once()
			})
			It("Akan Mengembalikan data user terbaru", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				udata, err := UserService.Update(ctx, entity.UpdateReq{Password: "satrio2323223", Email: "satrio44@gmail.com", Image: "gambar.jpg"}, file)
				Expect(err).Should(BeNil())
				Expect(udata.Email).To(Equal("satrio44@gmail.com"))
			})
		})

	})
	Context("User Delete", func() {
		When("Terjadi kesalahan pada database", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				err := UserService.Delete(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil menghapus akun", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan nil error", func() {
				err := UserService.Delete(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})
	Context("User Profile", func() {
		When("Id user tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("id not found")).Once()
			})
			It("Akan Mengembalikan Error dengan pesan 'id not found'", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("id not found"))
				Expect(user).To(BeNil())
			})
		})
		When("Server error", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan data user", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
				Expect(user).Should(BeNil())
			})
		})
		When("Id user ditemukan", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(&entity2.User{Email: "satrio@gmail.com"}, nil).Once()
			})
			It("Akan Mengembalikan data user", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(user.Email).To(Equal("satrio@gmail.com"))
			})
		})

	})

})
