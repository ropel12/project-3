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
	"github.com/ropel12/project-3/config/dependcy/container"
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
		Depend.Mds = container.NewMidtrans(&config.Config{Midtrans: config.MidtransConfig{ServerKey: "SB-Mid-server-TvgWB_Y9s81-rbMBH7zZ8BHW", ClientKey: "SB-Mid-client-nKsqvar5cn60u2Lv", Env: 1, ExpiryDuration: 1}})
		Mock = mocks.NewTransactionRepo(GinkgoT())
		TrxService = trx.NewTransactionService(Mock, Depend)
	})
	Context("Create Cart", func() {
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
			})
		})
	})

	Context("Checkout", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Eror dengan pesan 'Invalid or missing request body'", func() {
				res, err := TrxService.CreateTransaction(ctx, entity.ReqCheckout{})
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Invalid or missing request body"))
			})
		})
		When("Kouta yang dibeli melebihi jumlah kouta yang tersedia", func() {
			BeforeEach(func() {
				Mock.On("CheckQuota", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("You have exceeded the quota"))
			})
			It("Akan Mengembalikan Eror dengan pesan 'You have exceeded the quota'", func() {
				itemsdetails := []entity.ItemDetails{}
				itemsdetails = append(itemsdetails, entity.ItemDetails{Name: "Vip", Price: 200, SubTotal: 1000, Qty: 99})
				req := entity.ReqCheckout{
					EventId:     1,
					PaymentType: "bca",
					ItemDetails: itemsdetails,
				}
				res, err := TrxService.CreateTransaction(ctx, req)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("You have exceeded the quota"))
			})
		})
		When("Req body midtrans tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("CheckQuota", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				Mock.On("GetDetailUserById", mock.Anything, mock.Anything).Return(&entity2.User{Name: "satrio"})
			})
			It("Akan Mengembalikan Eror dengan pesan 'Missing or Invalid Request Body'", func() {
				itemsdetails := []entity.ItemDetails{}
				itemsdetails = append(itemsdetails, entity.ItemDetails{Name: "Vip", Price: 200, SubTotal: 1000})
				req := entity.ReqCheckout{
					EventId:     1,
					PaymentType: "bca",
					ItemDetails: itemsdetails,
				}
				res, err := TrxService.CreateTransaction(ctx, req)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Invalid Request Body"))
			})
		})
		When("Terdapat Kesalahan query database pada saat menyimpan data transaksi", func() {
			itemsdetails := []entity.ItemDetails{}
			req := entity.ReqCheckout{}
			BeforeEach(func() {
				itemsdetails = append(itemsdetails, entity.ItemDetails{Name: "Vip", Price: 1000, Qty: 1, SubTotal: 1000})
				req = entity.ReqCheckout{
					EventId:     1,
					PaymentType: "bca",
					ItemDetails: itemsdetails,
				}
				Mock.On("CheckQuota", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				Mock.On("GetDetailUserById", mock.Anything, mock.Anything).Return(&entity2.User{Name: "satrio"})
				Mock.On("CreateTransaction", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error"))
			})
			It("Akan Mengembalikan Eror dengan pesan 'Internal server error'", func() {
				res, err := TrxService.CreateTransaction(ctx, req)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil membuat transaksi dengan total transaksi sama dengan 0", func() {
			itemsdetails := []entity.ItemDetails{}
			req := entity.ReqCheckout{}
			BeforeEach(func() {
				itemsdetails = append(itemsdetails, entity.ItemDetails{Name: "regular", Price: 0, Qty: 1, SubTotal: 0})
				req = entity.ReqCheckout{
					EventId:     1,
					PaymentType: "free",
					ItemDetails: itemsdetails,
				}
				Mock.On("CheckQuota", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				Mock.On("GetDetailUserById", mock.Anything, mock.Anything).Return(&entity2.User{Name: "satrio"})
				Mock.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)
			})
			It("Akan Mengembalikan data Transaksi", func() {
				res, err := TrxService.CreateTransaction(ctx, req)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())
			})
		})
		When("Berhasil membuat transaksi dengan total transaksi lebih dari 0", func() {
			itemsdetails := []entity.ItemDetails{}
			req := entity.ReqCheckout{}
			BeforeEach(func() {
				itemsdetails = append(itemsdetails, entity.ItemDetails{Name: "Vip", Price: 1000, Qty: 1, SubTotal: 1000})
				req = entity.ReqCheckout{
					EventId:     1,
					PaymentType: "indomaret",
					ItemDetails: itemsdetails,
				}
				Mock.On("CheckQuota", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				Mock.On("GetDetailUserById", mock.Anything, mock.Anything).Return(&entity2.User{Name: "satrio"})
				Mock.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)
			})
			It("Akan Mengembalikan data Transaksi", func() {
				res, err := TrxService.CreateTransaction(ctx, req)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())
			})
		})

	})

	Context("Notification payment", func() {
		When("Terdapat kesalahan query db pada saat menyimpan data transaksi", func() {
			BeforeEach(func() {
				Mock.On("UpdateStatusTrasansaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Server Error"))
			})
			It("Akan Mengembalikan Error dengan pesan 'Internal Server Error'", func() {
				err := TrxService.UpdateStatus(ctx, "success", "INV-123323232")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})

		When("Sukses mengubah data transaksi", func() {
			BeforeEach(func() {
				Mock.On("UpdateStatusTrasansaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			})
			It("Akan Mengembalikan Error dengan nilai nil", func() {
				err := TrxService.UpdateStatus(ctx, "success", "INV-123323232")
				Expect(err).Should(BeNil())
			})
		})
	})
	Context("Detail Transaction", func() {
		When("Tidak terdapat data pada invoice dan user id yang di inputkan", func() {
			BeforeEach(func() {
				Mock.On("GetByInvoice", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found"))
			})
			It("Akan Mengembalikan Error dengan pesan 'Data Not Found'", func() {
				res, err := TrxService.GetDetail(ctx, "INV-123323232", 99)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Data Not Found"))
			})
		})
		When("Terdapat data pada invoice dan user id yang di inputkan", func() {
			BeforeEach(func() {
				TransactionItems := []entity2.TransactionItems{}
				TransactionItems = append(TransactionItems, entity2.TransactionItems{Qty: 1, Price: 1000})
				Transaction := entity2.Transaction{
					Status:           "paid",
					TransactionItems: TransactionItems,
				}
				Mock.On("GetByInvoice", mock.Anything, mock.Anything, mock.Anything).Return(&Transaction, nil)
			})
			It("Akan Mengembalikan data detail transaksi", func() {
				res, err := TrxService.GetDetail(ctx, "INV-123323232", 1)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())
			})
		})
	})

	Context("Get History By Uid", func() {
		When("Tidak Terdapat data pada user id yang di inputkan", func() {
			BeforeEach(func() {
				Mock.On("GetHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("Data not found"))
			})
			It("Akan Mengembalikan error dengan pesan 'Data not found'", func() {
				res, err := TrxService.GetHistoryByuid(ctx, 99, 1, 20)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Data not found"))
			})
		})
		When("Terdapat kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("GetHistory", mock.Anything, 99, mock.Anything, mock.Anything).Return(nil, 0, errors.New("Internal server error"))
			})
			It("Akan Mengembalikan error dengan pesan 'Internal server error'", func() {
				res, err := TrxService.GetHistoryByuid(ctx, 99, 1, 10)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Terdapat data pada user id yang di inputkan", func() {
			BeforeEach(func() {
				Event := entity2.Event{Name: "Dota2"}
				trx := entity2.Transaction{Invoice: "INV-12121212121", Event: Event}
				datatrx := []entity2.Transaction{}
				datatrx = append(datatrx, trx)
				Mock.On("GetHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(datatrx, 0, nil)
			})
			It("Akan Mengembalikan data transaksi user", func() {
				res, err := TrxService.GetHistoryByuid(ctx, 1, 10, 20)
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
			})
		})
	})

	Context("Get Transaction By Status", func() {
		When("Status bukan paid atau pending", func() {
			It("Akan Mengembalikan error dengan pesan 'Data not found'", func() {
				res, err := TrxService.GetByStatus(ctx, 1, "Ngasal")
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Data not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("GetByStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Internal server error"))
			})
			It("Akan Mengembalikan error dengan pesan 'Internal server error", func() {
				res, err := TrxService.GetByStatus(ctx, 1, "paid")
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Terdapat data dengan status yang di inputkan", func() {
			BeforeEach(func() {
				trx := entity2.Transaction{Invoice: "INV-12121212121", Event: entity2.Event{Name: "Dota22"}}
				datatrx := []entity2.Transaction{}
				datatrx = append(datatrx, trx)
				Mock.On("GetByStatus", mock.Anything, mock.Anything, mock.Anything).Return(datatrx, nil)
			})
			It("Akan Mengembalikan data transaksi", func() {
				res, err := TrxService.GetByStatus(ctx, 1, "paid")
				Expect(err).Should(BeNil())
				Expect(res.Data).ShouldNot(BeNil())
			})
		})
	})

	Context("Get Ticket Transaction", func() {
		When("Tidak terdapat data pada invoice yang dimasukan", func() {
			BeforeEach(func() {
				Mock.On("GetTicketByInvoice", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data not found"))
			})
			It("Akan Mengembalikan error dengan pesan 'Data not found", func() {
				res, err := TrxService.GetTickets(ctx, "INV-999999999999", 99)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Data not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("GetTicketByInvoice", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Internal server error"))
			})
			It("Akan Mengembalikan error dengan pesan 'Internal server error", func() {
				res, err := TrxService.GetTickets(ctx, "INV-1231113", 1)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
				Expect(err.Error()).To(Equal("Internal server error"))
			})
		})
		When("Terdapat data pada invoice yang dimasukan", func() {
			BeforeEach(func() {
				event := entity2.Event{Name: "Dota2", Location: "Senayan"}
				trxitems := []entity2.TransactionItems{entity2.TransactionItems{Qty: 1, Type: entity2.Type{Name: "regular"}}}
				res := entity2.Transaction{Event: event, TransactionItems: trxitems}
				Mock.On("GetTicketByInvoice", mock.Anything, mock.Anything, mock.Anything).Return(&res, nil)
			})
			It("Akan Mengembalikan data ticket dari invoice tersebut", func() {
				res, err := TrxService.GetTickets(ctx, "INV-1231113", 1)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())
			})
		})
	})
})
