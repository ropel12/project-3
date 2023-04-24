package service_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	entity2 "github.com/ropel12/project-3/app/entities"
	entity "github.com/ropel12/project-3/app/features/transaction"
	mocks "github.com/ropel12/project-3/app/features/transaction/mocks/repository"
	trx "github.com/ropel12/project-3/app/features/transaction/service"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/sirupsen/logrus"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("user", func() {
	var Mock *mocks.TransactionRepo
	var TrxService trx.TransactionService
	var Depend dependcy.Depend
	var ctx context.Context
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		Mock = mocks.NewTransactionRepo(GinkgoT())
		TrxService = trx.NewTransactionService(Mock, Depend)
	})
	Context("User Login", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Missing or Invalid Request Body'", func() {
				err := TrxService.CreateCart(ctx, entity.ReqCart{})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Missing or Invalid Request Body"))
			})
		})
		When("Kesalahan Query database", func() {
			data := entity2.Carts{
				UserID: 1,
				TypeID: 1,
				Qty:    1,
			}
			BeforeEach(func() {
				Mock.On("Create", mock.Anything, data).Return(errors.New("Internal server error"))
			})
			It("Akan Mengembalikan Eror dengan pesan 'Internal server error'", func() {
				err := TrxService.CreateCart(ctx, entity.ReqCart{TypeID: 1, UID: 1})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Sukses melakukan operasi", func() {
			data := entity2.Carts{
				UserID: 1,
				TypeID: 1,
				Qty:    1,
			}
			BeforeEach(func() {
				Mock.On("Create", mock.Anything, data).Return(nil)
			})
			It("Akan Mengembalikan Eror dengan nilai nil", func() {
				err := TrxService.CreateCart(ctx, entity.ReqCart{TypeID: 1, UID: 1})
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("Get Cart Byuid", func() {
		When("Data cart dengan uid sekarang tidak ada", func() {
			BeforeEach(func() {
				Mock.On("GetCart", mock.Anything, 1).Return(nil, errors.New("Data not found"))
			})
			It("Akan Mengembalikan error dengan pesan 'Data not found'", func() {
				data, err := TrxService.GetCart(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(data).Should(BeNil())
				Expect(err.Error()).To(Equal("Data not found"))
			})
		})

		When("Kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("GetCart", mock.Anything, 1).Return(nil, errors.New("Internal Server Error"))
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error", func() {
				data, err := TrxService.GetCart(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(data).Should(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Terdapat data pada uid sekarang", func() {
			BeforeEach(func() {
				data := []entity2.Carts{}
				data = append(data, entity2.Carts{UserID: 1, TypeID: 1})
				Mock.On("GetCart", mock.Anything, 1).Return(data, nil)
			})
			It("Akan Mengembalikan data cart", func() {
				data, err := TrxService.GetCart(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
				Expect(data[0].TypeID).To(Equal(1))
			})
		})
	})
})
