package model

type OfficialAccount struct {
	Id              int    `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT" json:"id"`
	NickName        string `gorm:"column:nick_name;type:varchar(100)" json:"nick_name"`
	Appid           string `gorm:"column:appid;type:varchar(100)" json:"appid"`
	HeadImg         string `gorm:"column:head_img;type:varchar(200)" json:"head_img"`
	ServiceTypeInfo int    `gorm:"column:service_type_info;type:tinyint(1)" json:"service_type_info"`
	VerifyTypeInfo  int    `gorm:"column:verify_type_info;type:tinyint(1)" json:"verify_type_info"`
	UserName        string `gorm:"column:user_name;type:varchar(100)" json:"user_name"`
	PrincipalName   string `gorm:"column:principal_name;type:varchar(100)" json:"principal_name"`
	Alias           string `gorm:"column:alias;type:varchar(50)" json:"alias"`
	Signature       string `gorm:"column:signature;type:varchar(200)" json:"signature"`
	OpenStore       string `gorm:"column:open_store;type:tinyint(1)" json:"open_store"`
	OpenScan        string `gorm:"column:open_scan;type:tinyint(1)" json:"open_scan"`
	OpenPay         string `gorm:"column:open_pay;type:tinyint(1)" json:"open_pay"`
	OpenCard        string `gorm:"column:open_card;type:tinyint(1)" json:"open_card"`
	OpenShake       string `gorm:"column:open_shake;type:tinyint(1)" json:"open_shake"`
	QrcodeUrl       string `gorm:"column:qrcode_url;type:varchar(500)" json:"qrcode_url"`
}

func (m *OfficialAccount) TableName() string {
	return "official_account"
}
