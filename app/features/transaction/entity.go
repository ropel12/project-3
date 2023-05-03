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
		Name     string `json:"type_name,omitempty" validate:"required"`
		Price    int    `json:"type_price,omitempty" validate:"required"`
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
		Id           int    `json:"id,omitempty"`
		Name         string `json:"name,omitempty"`
		Date         string `json:"date,omitempty"`
		Location     string `json:"location,omitempty"`
		EndDate      string `json:"end_date,omitempty"`
		HostedBy     string `json:"hosted_by,omitempty"`
		Image        string `json:"image,omitempty"`
		Participants int    `json:"participants,omitempty"`
	}
	Transaction struct {
		Total         int64         `json:"total,omitempty"`
		Date          string        `json:"date,omitempty"`
		Expire        string        `json:"expire,omitempty"`
		PaymentMethod string        `json:"payment_method,omitempty"`
		Invoice       string        `json:"invoice,omitempty"`
		EventName     string        `json:"event_name,omitempty"`
		Status        string        `json:"status,omitempty"`
		PaymentCode   string        `json:"payment_code,omitempty"`
		ItemDetails   []ItemDetails `json:"item_details,omitempty"`
	}
	TicketTransaction struct {
		TicketType string `json:"ticket_type"`
		EventName  string `json:"event_name"`
		Location   string `json:"location"`
		Date       string `json:"date"`
		HostedBy   string `json:"hosted_by"`
		Qty        int    `json:"-"`
	}

	Response struct {
		Csrf  string `json:"csrf,omitempty"`
		Total int    `json:"total,omitempty"`
		Data  any    `json:"data"`
	}
)
