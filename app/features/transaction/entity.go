package transaction

type (
	ReqCart struct {
		UID    int
		TypeID int `json:"type_id" validate:"required"`
	}
	Cart struct {
		EventId   int    `json:"event_id" `
		TypeID    int    `json:"type_id"`
		TypeName  string `json:"type_name"`
		TypePrice int    `json:"type_price"`
		Qty       int    `json:"qty"`
		Subtotal  int    `json:"sub_total"`
	}
	ItemDetails struct {
		Name     string `json:"name" validate:"required"`
		Price    int    `json:"price" validate:"required"`
		TypeId   int    `json:"type_id" validate:"required"`
		SubTotal int    `json:"sub_total" validate:"required"`
		Qty      int    `json:"qty" validate:"required"`
	}
	ReqCheckout struct {
		EventId     int           `json:"event_id" validate:"required" `
		PaymentType string        `json:"payment_type" validate:"required"`
		ItemDetails []ItemDetails `json:"items_detail" validate:"required"`
		UserId      int
	}
)
