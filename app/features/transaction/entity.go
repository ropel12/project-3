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
	EventTransaction struct {
		Id           int    `json:"id"`
		Date         string `json:"date"`
		Location     string `json:"location"`
		EndDate      string `json:"end_date"`
		HostedBy     string `json:"hosted_by"`
		Image        string `json:"image"`
		Participants int    `json:"participants"`
	}
	Transaction struct {
		Total         int64         `json:"total,omitempty"`
		Date          string        `json:"date,omitempty"`
		Expire        string        `json:"expire,omitempty"`
		PaymentMethod string        `json:"payment_method,omitempty"`
		Status        string        `json:"status,omitempty"`
		PaymentCode   string        `json:"payment_code,omitempty"`
		ItemDetails   []ItemDetails `json:"item_details,omitempty"`
	}

	Response struct {
		Data any `json:"data"`
	}
)
