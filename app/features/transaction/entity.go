package transaction

type (
	ReqCart struct {
		UID    int
		TypeID int `json:"type_id" validate:"required"`
	}
)
