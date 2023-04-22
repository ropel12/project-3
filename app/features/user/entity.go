package user

type (
	LoginReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	RegisterReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Address  string `json:"address" validate:"required"`
	}
	UpdateReq struct {
		Id       int
		Email    string `form:"email"`
		Password string `form:"password"`
		Name     string `form:"name" `
		Address  string `form:"address"`
		Image    string `form:"image"`
	}
)
