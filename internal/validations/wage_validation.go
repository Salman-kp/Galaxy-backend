package validations

import "errors"

type UpdateWageRequest struct {
	TAAmount    int64 `json:"ta_amount"`
	BonusAmount int64 `json:"bonus_amount"`
	FineAmount  int64 `json:"fine_amount"`
}

func (r *UpdateWageRequest) Validate() error {
	if r.TAAmount < 0 {
		return errors.New("TA amount cannot be negative")
	}
	if r.BonusAmount < 0 {
		return errors.New("bonus amount cannot be negative")
	}
	if r.FineAmount < 0 {
		return errors.New("fine amount cannot be negative")
	}
	return nil
}