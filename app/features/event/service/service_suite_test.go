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
	"github.com/ropel12/project-3/config/dependcy/container"
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
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		Depend.Rds = redis.NewClient(&redis.Options{})
		Mock = mocks.NewEventRepo(GinkgoT())
		validation, _ := container.NewValidation()
		Depend.Validation = validation
		EventService = event.NewEventService(Mock, Depend)
	})
	Context("Create Event", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Invalid or missing request body'", func() {
				var file multipart.File
				id, err := EventService.Create(ctx, entity.ReqCreate{}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
				Expect(id).To(Equal(0))
			})
		})
		When("Terdapat salah satu request body yang tidak diisi", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Invalid or missing request body'", func() {
				var file multipart.File
				id, err := EventService.Create(ctx, entity.ReqCreate{Name: "Dota2"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
				Expect(id).To(Equal(0))
			})
		})
		When("Terdapat request gambar yang tidak sesuai dengan format gambar", func() {
			It("Akan Mengembalikan Eror dengan pesan 'File type not allowed'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreate{
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2023-05-19T14:50:32",
					Quota:     10,
					Image:     "image.js",
					Rtype:     `[{"id":55,"type_name":"vvvip","price":200000}]`,
				}
				id, err := EventService.Create(ctx, req, file)
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
				req := entity.ReqCreate{
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2023-05-19T14:50:32",
					Quota:     10,
					Image:     "image.jpg",
					Rtype:     `[{"id":55,"type_name":"vvvip","price":200000}]`,
				}
				id, err := EventService.Create(ctx, req, file)
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
				req := entity.ReqCreate{
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2023-05-19T14:50:32",
					Quota:     10,
					Image:     "image.jpg",
					Rtype:     `[{"id":55,"type_name":"vvvip","price":200000}]`,
				}
				file = os.NewFile(uintptr(2), "2")
				id, err := EventService.Create(ctx, req, file)
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
				var Types []entity2.Type
				Comments = append(Comments, entity2.UserComments{UserID: 1})
				Users = append(Users, entity2.User{Name: "satrio"})
				Types = append(Types, entity2.Type{Name: "dota2", Price: 9000})
				res := &entity2.Event{Name: "Dota 2", Users: Users, UserComments: Comments, Types: Types}
				Mock.On("GetById", mock.Anything, 1).Return(res, nil).Once()
			})
			It("Akan Mengembalikan data event", func() {
				res, err := EventService.Detail(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
			})
		})

	})

	Context("Update Event", func() {
		When("Request body tidak valid atau tidak ada", func() {
			It("Akan Mengembalikan error dengan pesan 'Invalid or missing request body'", func() {
				var file multipart.File
				id, err := EventService.Update(ctx, entity.ReqUpdate{}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
			})
		})
		When("Req body image bukan merupakan gambar", func() {
			It("Akan Mengembalikan error dengna pesan 'File type not allowed'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdate{
					Id:        1,
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2022-05-06 15:04:05",
					Quota:     10,
					Image:     "image.js",
					Rtype:     "[{\"id\":55,\"type_name\":\"vvvip\",\"price\":200000}]",
				}
				id, err := EventService.Update(ctx, req, file)
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Terdapat kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Internal Server Error'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdate{
					Id:        1,
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2022-05-06 15:04:05",
					Quota:     10,
					Image:     "image.jpg",
					Rtype:     `[{"id":55,"type_name":"vvvip","price":200000}]`,
				}
				id, err := EventService.Update(ctx, req, file)
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbarui data event", func() {
			BeforeEach(func() {
				res := entity2.Event{
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Detail:    "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2022-05-06 15:04:05",
					Quota:     10,
					Image:     "image.jpg",
				}
				res.ID = 1
				Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&res, nil).Once()
			})
			It("Akan Mengembalikan id event", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdate{
					Id:        1,
					Name:      "Linux",
					Location:  "JKT",
					Duration:  2.5,
					Details:   "Belajar Linux",
					HostedBy:  "Claudio delvin",
					StartDate: "2022-05-06 15:04:05",
					Quota:     10,
					Image:     "image.jpg",
					Rtype:     `[{"id":55,"type_name":"vvvip","price":200000}]`,
					Types:     []entity.TypeEvent{entity.TypeEvent{Id: 1, TypeName: "Dota", Price: 2000}},
				}
				id, err := EventService.Update(ctx, req, file)
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})
	})
	Context("Create comment", func() {
		When("Request body kosong atau tidak sesuai", func() {
			It("Akan Mengembalikan error dengna pesan 'Invalid or missing request body'", func() {

				id, err := EventService.CreateComment(ctx, entity.ReqCreateComment{})
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
			})
		})
		When("Terdapat kesalahan database", func() {
			BeforeEach(func() {
				Mock.On("CreateComment", mock.Anything, mock.Anything).Return(nil, errors.New("Internal server error")).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Internal server error'", func() {
				id, err := EventService.CreateComment(ctx, entity.ReqCreateComment{EventId: 1, Comment: "Eventnya bagus"})
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Berhasil membuat comment", func() {
			BeforeEach(func() {
				Mock.On("CreateComment", mock.Anything, mock.Anything).Return(&entity2.UserComments{EventID: 1}, nil).Once()
			})
			It("Akan Mengembalikan event id", func() {
				id, err := EventService.CreateComment(ctx, entity.ReqCreateComment{EventId: 1, Comment: "Eventnya bagus"})
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})
	})
	Context("Create Ticket", func() {
		When("Request body kosong atau tidak sesuai", func() {
			It("Akan Mengembalikan error dengna pesan 'Invalid or missing request body'", func() {
				id, err := EventService.CreateTicket(ctx, entity.ReqCreateTicket{})
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
			})
		})

		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("CreateTicket", mock.Anything, mock.Anything).Return(nil, errors.New("Internal server error")).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Internal server error'", func() {
				id, err := EventService.CreateTicket(ctx, entity.ReqCreateTicket{EventId: 1, TypeName: "VIp", Price: 2000})
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				typee := entity2.Type{EventID: 1}

				Mock.On("CreateTicket", mock.Anything, mock.Anything).Return(&typee, nil).Once()
			})
			It("Akan Mengembalikan error dengna pesan 'Internal server error'", func() {
				id, err := EventService.CreateTicket(ctx, entity.ReqCreateTicket{EventId: 1, TypeName: "VIp", Price: 2000})
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})
	})
})
