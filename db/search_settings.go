package db

import "time"

type SearchSetting struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	SearchOn  bool      `json:"searchOn"`
	AddNew    bool      `json:"addNew"`
	Amount    uint      `json:"amount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *SearchSetting) Get() error {
	err := DBConnect.Where("id = 1").First(s).Error
	return err
}

func (s *SearchSetting) Update() error {
	tx := DBConnect.Select("search_on", "add_new", "amount", "updated_at").Where("id = 1").Updates(s)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
