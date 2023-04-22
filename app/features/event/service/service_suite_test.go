package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				Mock.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("error database"))

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
				Mock.On("Create", mock.Anything, mock.Anything).Return(&id, nil)
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
})
