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
)
