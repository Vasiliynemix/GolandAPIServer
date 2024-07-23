package structures

// UserSetRequest model info
// @Description User setter
type UserSetRequest struct {
	TgID      int64   `json:"tg_id" format:"int64" example:"54321" binding:"required"` // This is TgID in Telegram
	UserName  *string `json:"user_name" format:"string" example:"username"`            // This is UserName in Telegram, may be nil
	LastName  *string `json:"last_name" format:"string" example:"lastname"`            // This is LastName in Telegram, may be nil
	FirstName *string `json:"first_name" format:"string" example:"firstname"`          // This is FirstName in Telegram, may be nil
}

// UserShow model info
// @Description Show User Info
type UserShow struct {
	TgID      int64   `json:"tg_id" format:"int64" example:"54321" binding:"required"`           // This is TgID in Telegram
	UserName  *string `json:"user_name" format:"string" example:"username" binding:"required"`   // This is UserName in Telegram, may be nil
	LastName  *string `json:"last_name" format:"string" example:"lastname" binding:"required"`   // This is LastName in Telegram, may be nil
	FirstName *string `json:"first_name" format:"string" example:"firstname" binding:"required"` // This is FirstName in Telegram, may be nil
}
