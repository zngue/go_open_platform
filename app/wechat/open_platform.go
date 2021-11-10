package wechat

import (
	"encoding/xml"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/openplatform"
	openConfig "github.com/silenceper/wechat/v2/openplatform/config"
	"github.com/silenceper/wechat/v2/util"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg"
	"github.com/zngue/go_open_platform/app/model"
	"time"
)

type OpenPlatform struct {
	config *openConfig.Config
}
type MessageEncrypt struct {
	AppId string `xml:"AppId"`
	Encrypt string `xml:"Encrypt"`
}
type ComponentVerifyTicketEncrypt struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}
type Authorizer struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

type NotifyThirdFasterRegister struct {
	Status int `xml:"status"`
	Msg string `xml:"msg"`
	Info struct{
		Name string `xml:"name"`
		Code  string `xml:"code"`
		CodeType int `xml:"code_type"`
		LegalPersonaWechat string `xml:"legal_persona_wechat"`
		LegalPersonaName string `xml:"legal_persona_name"`
		ComponentPhone string `xml:"component_phone"`
	} `xml:"info"`
}

type BaseMessageDecEncrypt struct {
	InfoType string `xml:"InfoType"`
	AppId string `xml:"AppId"`
	CreateTime int64 `xml:"CreateTime"`
	XmlByte  string
	ComponentVerifyTicketEncrypt
	NotifyThirdFasterRegister
	Authorizer
}

func  (o *BaseMessageDecEncrypt) Init()  {
	if o.InfoType=="component_verify_ticket" {
		verifyTime := viper.GetInt64("wechatOpenPlatform.VerifyTime")
		if verifyTime+3600*9<o.CreateTime {
			viper.Set("wechatOpenPlatform.VerifyTicket",o.ComponentVerifyTicket)
			viper.Set("wechatOpenPlatform.VerifyTime",o.CreateTime)
		}
	}
	if len(o.XmlByte)>0 {
		go o.SaveMessage()
	}
}
func  (o *BaseMessageDecEncrypt) SaveMessage()  {
	pkg.MysqlConn.Model(&model.MessageContent{}).Create(&model.MessageContent{
		Content: o.XmlByte,
		CreateTime: time.Now().Unix(),
	})
}
type IOpenPlatform interface {
	Platform(isToken bool) *openplatform.OpenPlatform
	DecryptMsg(message []byte) (*BaseMessageDecEncrypt,error)
}

func (o *OpenPlatform) DecryptMsg(message []byte) (*BaseMessageDecEncrypt,error) {
	mData := new(MessageEncrypt)
	if err := xml.Unmarshal(message, mData); err != nil {
		return nil, err
	}
	var data BaseMessageDecEncrypt
	_, xmlBytes, xmlErr := util.DecryptMsg(o.config.AppID, mData.Encrypt, o.config.EncodingAESKey)
	if xmlErr != nil {
		return nil, xmlErr
	}
	data.XmlByte=string(xmlBytes)
	if err := xml.Unmarshal(xmlBytes, &data); err != nil {
		return nil, err
	}
	if &data!=nil {
		data.Init()
	}
	return &data, nil
}
func (o *OpenPlatform) Platform(isToken bool) *openplatform.OpenPlatform {
	newWechat := wechat.NewWechat()
	platform := newWechat.GetOpenPlatform(o.config)
	if isToken {
		verifyTicket := viper.GetString("wechatOpenPlatform.VerifyTicket")
		token, err := platform.GetComponentAccessToken()
		if err!=nil || token=="" {
			platform.SetComponentAccessToken(verifyTicket)
		}
	}
	return platform
}

func NewOpenPlatform() IOpenPlatform {
	memory := cache.NewMemory()
	config:=&openConfig.Config{
		AppID:          viper.GetString("wechatOpenPlatform.Appid"),
		AppSecret:      viper.GetString("wechatOpenPlatform.AppSecret"),
		Token:          viper.GetString("wechatOpenPlatform.Token"),
		EncodingAESKey: viper.GetString("wechatOpenPlatform.EncodingAESKey"),
		Cache:          memory,
	}
	platform := new(OpenPlatform)


	platform.config=config
	return platform
}

