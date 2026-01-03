package captain

type BookingDTO struct {
	ID          uint   `json:"id"`
	EventID     uint   `json:"event_id"`
	Status      string `json:"status"`
	Role        string `json:"role"`
	BaseAmount  int64  `json:"base_amount"`
	ExtraAmount int64  `json:"extra_amount"`
	TAAmount    int64  `json:"ta_amount"`
	BonusAmount int64  `json:"bonus_amount"`
	FineAmount  int64  `json:"fine_amount"`
	TotalAmount int64  `json:"total_amount"`
	CreatedAt   string `json:"created_at"`
}

type AttendanceRowResponse struct {
	BookingID  uint   `json:"booking_id"`
	UserID     uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	Role       string `json:"role"`
	Status     string `json:"status"`

	BaseAmount  int64 `json:"base_amount"`
	ExtraAmount int64 `json:"extra_amount"`
	TAAmount    int64 `json:"ta_amount"`
	BonusAmount int64 `json:"bonus_amount"`
	FineAmount  int64 `json:"fine_amount"`
	TotalAmount int64 `json:"total_amount"`
}
