package wechat

import "github.com/spf13/viper"

type AuthLinkRequest struct {
	CallbackUrl string `form:"callback_url" json:"callback_url"`
	AuthType    int    `json:"auth_type" form:"auth_type"`
	BizAppID    string `form:"auth_type"`
	IsMobile    bool   `form:"is_mobile"`
}
type AuthLinkRsp struct {
	Code string `json:"code"`
	Link string `json:"link"`
}

func (a *AuthLinkRequest) Init() {
	if a.AuthType == 0 {
		a.AuthType = 3
	}
}

type MessageEncrypt struct {
	AppId   string `xml:"AppId"`
	Encrypt string `xml:"Encrypt"`
}
type ComponentVerifyTicketEncrypt struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}
type Authorizer struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}
type NotifyThirdFasterRegister struct {
	Status int    `xml:"status"`
	Msg    string `xml:"msg"`
	Info   struct {
		Name               string `xml:"name"`
		Code               string `xml:"code"`
		CodeType           int    `xml:"code_type"`
		LegalPersonaWechat string `xml:"legal_persona_wechat"`
		LegalPersonaName   string `xml:"legal_persona_name"`
		ComponentPhone     string `xml:"component_phone"`
	} `xml:"info"`
}

type BaseMessageDecEncrypt struct {
	InfoType   string `xml:"InfoType"`
	AppId      string `xml:"AppId"`
	CreateTime int64  `xml:"CreateTime"`
	XmlByte    string
	OriginXml  string
	ComponentVerifyTicketEncrypt
	NotifyThirdFasterRegister
	Authorizer
}

func (o *BaseMessageDecEncrypt) Init() {
	if o.InfoType == "component_verify_ticket" {
		verifyTime := viper.GetInt64("wechatOpenPlatform.VerifyTime")
		if verifyTime+3600*9 < o.CreateTime {
			viper.Set("wechatOpenPlatform.VerifyTicket", o.ComponentVerifyTicket)
			viper.Set("wechatOpenPlatform.VerifyTime", o.CreateTime)
		}
	}
	if len(o.XmlByte) > 0 {
		go o.SaveMessage()
	}
}
