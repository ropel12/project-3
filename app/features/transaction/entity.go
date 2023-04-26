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
		Name     string `json:"name,omitempty" validate:"required"`
		Price    int    `json:"price,omitempty" validate:"required"`
		TypeId   int    `json:"type_id,omitempty" validate:"required"`
		SubTotal int    `json:"sub_total,omitempty" validate:"required"`
		Qty      int    `json:"qty,omitempty" validate:"required"`
	}
	ReqCheckout struct {
		EventId     int           `json:"event_id" validate:"required" `
		PaymentType string        `json:"payment_type" validate:"required"`
		ItemDetails []ItemDetails `json:"items_detail" validate:"required"`
		UserId      int
	}
	Transaction struct {
		Total         int64         `json:"total"`
		Date          string        `json:"date"`
		Expire        string        `json:"expire"`
		PaymentMethod string        `json:"payment_method"`
		Status        string        `json:"status"`
		PaymentCode   string        `json:"payment_code"`
		ItemDetails   []ItemDetails `json:"item_details"`
	}
	Response struct {
		Data any `json:"data"`
	}
)
