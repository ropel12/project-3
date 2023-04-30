package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	entity2 "github.com/ropel12/project-3/app/entities"
	entity "github.com/ropel12/project-3/app/features/event"
	mocks "github.com/ropel12/project-3/app/features/event/mocks/repository"
	event "github.com/ropel12/project-3/app/features/event/service"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("event", func() {
	var Mock *mocks.EventRepo
	var EventService event.EventService
	var Depend dependcy.Depend
	var ctx context.Context
	var FailedReq entity.ReqCreate
	var SuccessReq entity.ReqCreate
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		Depend.Rds = redis.NewClient(&redis.Options{})
		Mock = mocks.NewEventRepo(GinkgoT())
		EventService = event.NewEventService(Mock, Depend)
		FailedReq = entity.ReqCreate{Image: "gambar.js", Name: "tes", StartDate: "2006-01-02 15:04:05", Duration: 2.5, Details: "tes", Location: "jakarta", Quota: 100, Rtype: `[{"name":"vip","price":2000},{"name":"regular","price":1000}]`, HostedBy: "ropel"}
		SuccessReq = entity.ReqCreate{Image: "gambar.jpg", Name: "tes", StartDate: "2006-01-02 15:04:05", Duration: 2.5, Details: "tes", Location: "jakarta", Quota: 100, Rtype: `[{"name":"vip","price":2000},{"name":"regular","price":1000}]`, HostedBy: "ropel"}
	})
	Context("Create Event", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Invalid and missing request body'", func() {
				var file multipart.File
				id, err := EventService.Create(ctx, entity.ReqCreate{}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Invalid and missing request body"))
				Expect(id).To(Equal(0))
			})
		})
		When("Terdapat salah satu request body yang tidak diisi", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Invalid and missing request body'", func() {
				var file multipart.File
				id, err := EventService.Create(ctx, entity.ReqCreate{Name: "Dota2"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Invalid and missing request body"))
				Expect(id).To(Equal(0))
			})
		})
		When("Terdapat request gambar yang tidak sesuai dengan format gambar", func() {
			It("Akan Mengembalikan Eror dengan pesan 'File type not allowed'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				id, err := EventService.Create(ctx, FailedReq, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
				Expect(id).To(Equal(0))
			})

		})

		When("Terjadi kesalahn dalam database", func() {
			BeforeEach(func() {
				Mock.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("error database")).Once()

			})
			It("Akan Mengembalikan Eror dengan pesan 'error database'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				id, err := EventService.Create(ctx, SuccessReq, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("error database"))
				Expect(id).To(Equal(0))
			})

		})

		When("Sukses menambahkan event", func() {
			BeforeEach(func() {
				id := 1
				Mock.On("Create", mock.Anything, mock.Anything).Return(&id, nil).Once()
			})
			It("Akan Mengembalikan id event", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				id, err := EventService.Create(ctx, SuccessReq, file)
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})

		})

	})

	Context("Get My Event", func() {
		When("Terjadi kesalahan pada database", func() {
			uid := 1
			limit := 5
			offset := 0
			BeforeEach(func() {
				Mock.On("GetByUid", mock.Anything, mock.Anything, uid, limit, offset).Return(nil, 0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				_, err := EventService.MyEvent(ctx, uid, limit, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})

		When("User id tidak ditemukan", func() {
			limit := 5
			offset := 0
			BeforeEach(func() {
				Mock.On("GetByUid", mock.Anything, mock.Anything, mock.Anything, limit, offset).Return(nil, 0, errors.New("data not found")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'data not found'", func() {
				_, err := EventService.MyEvent(ctx, 99, limit, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("data not found"))
			})
		})

		When("Data Event User Ditemukan", func() {
			limit := 5
			offset := 0
			BeforeEach(func() {
				res := []*entity2.Event{}
				res = append(res, &entity2.Event{Name: "Dota 2"})
				Mock.On("GetByUid", mock.Anything, mock.Anything, mock.Anything, limit, offset).Return(res, 10, nil).Once()
			})
			It("Akan Mengembalikan data event yang dimiliki oleh user", func() {
				res, err := EventService.MyEvent(ctx, 1, limit, 1)
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
				Expect(res.Limit).To(Equal(5))
			})
		})
	})

	Context("Delete Event", func() {
		When("Event id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("id not found")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'id not found'", func() {
				err := EventService.Delete(ctx, 9, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("id not found"))
			})
		})

		When("Userid berbeda dengan uid pemilik event", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Cannot delete event")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Cannot delete event'", func() {
				err := EventService.Delete(ctx, 1, 5)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Cannot delete event"))
			})
		})

		When("Kesalahan pada database", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("internal server error"))
			})
			It("Akan Mengembalikan error dengan pesan 'internal server error'", func() {
				err := EventService.Delete(ctx, 1, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("internal server error"))
			})
		})

		When("Sukses menghapus event", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan err dengan nilai nil", func() {
				err := EventService.Delete(ctx, 1, 1)
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("Get All Event", func() {
		When("Terjadi kesalahan pada database", func() {
			limit := 5
			offset := 0
			BeforeEach(func() {
				Mock.On("GetAll", mock.Anything, mock.Anything, limit, offset).Return(nil, 0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				_, err := EventService.GetAll(ctx, limit, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})

		When("Data event tidak ditemukan", func() {
			limit := 5
			offset := 0
			BeforeEach(func() {
				Mock.On("GetAll", mock.Anything, mock.Anything, limit, offset).Return(nil, 0, errors.New("data not found")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'data not found'", func() {
				_, err := EventService.GetAll(ctx, limit, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("data not found"))
			})
		})

		When("Data Event Ditemukan", func() {
			limit := 5
			offset := 0
			BeforeEach(func() {
				res := []*entity2.Event{}
				res = append(res, &entity2.Event{Name: "Dota 2"})
				Mock.On("GetAll", mock.Anything, mock.Anything, limit, offset).Return(res, 10, nil).Once()
			})
			It("Akan Mengembalikan data event", func() {
				res, err := EventService.GetAll(ctx, limit, 1)
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
				Expect(res.Limit).To(Equal(5))
			})
		})
	})
	Context("Detail Event", func() {
		When("Tidak terdapat data pada id yang di masukan", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, 99).Return(nil, errors.New("Data not found")).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Data not found'", func() {
				_, err := EventService.Detail(ctx, 99)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Data not found"))
			})
		})

		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, 1).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Internal Server Error'", func() {
				_, err := EventService.Detail(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Terdapat data pada id yang diinputkan", func() {
			BeforeEach(func() {
				var Users []entity2.User
				var Comments []entity2.UserComments
				Comments = append(Comments, entity2.UserComments{UserID: 1})
				Users = append(Users, entity2.User{Name: "satrio"})
				res := &entity2.Event{Name: "Dota 2", Users: Users, UserComments: Comments}
				Mock.On("GetById", mock.Anything, 1).Return(res, nil).Once()
			})
			It("Akan Mengembalikan data event", func() {
				res, err := EventService.Detail(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
			})
		})

	})
})
