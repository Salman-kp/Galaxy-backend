package admin


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

type EventWageSummary struct {
	TotalWorkers       int   `json:"total_workers"`
	TotalBaseAmount    int64 `json:"total_base_amount"`
	TotalExtraAmount   int64 `json:"total_extra_amount"`
	TotalTAAmount      int64 `json:"total_ta_amount"`
	TotalBonusAmount   int64 `json:"total_bonus_amount"`
	TotalFineAmount    int64 `json:"total_fine_amount"`
	GrandTotalAmount   int64 `json:"grand_total_amount"`
}