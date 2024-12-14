package model

type Organization struct {
	ID        string `json:"id" gorm:"type:char(19);primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(64)"`
	CreatedAt string `json:"createdAt"`
}

type Member struct {
	AccountID string `json:"accountId" gorm:"type:char(19);primaryKey;index"`
	RoleID    string `json:"roleId" gorm:"type:char(19);primaryKey;index"`
}
