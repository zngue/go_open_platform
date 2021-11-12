package wechat

import (
	"encoding/xml"
	"errors"
	"fmt"
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
	config   *openConfig.Config
	platform *openplatform.OpenPlatform
}

func (o *BaseMessageDecEncrypt) SaveMessage() {
	pkg.MysqlConn.Model(&model.MessageContent{}).Create(&model.MessageContent{
		Content:    o.XmlByte,
		CreateTime: time.Now().Unix(),
	})
}

type IOpenPlatform interface {
	Platform(isToken bool) *openplatform.OpenPlatform
	DecryptMsg(message []byte) (*BaseMessageDecEncrypt, error)
	GetLinkByCode(code string) (string, error)
	AuthLink(req *AuthLinkRequest) (authLin *AuthLinkRsp, err error)
}

func (o *OpenPlatform) GetLinkByCode(code string) (string, error) {
	link := o.platform.Cache.Get(code)
	if linkUrl, ok := link.(string); ok {
		return linkUrl, nil
	} else {
		return "", errors.New("not exit code link")
	}
}
func (o *OpenPlatform) AuthLink(req *AuthLinkRequest) (authLin *AuthLinkRsp, err error) {
	req.Init()
	var link string
	if req.IsMobile {
		link, err = o.platform.GetBindComponentURL(req.CallbackUrl, req.AuthType, req.BizAppID)
	} else {
		link, err = o.platform.GetComponentLoginPage(req.CallbackUrl, req.AuthType, req.BizAppID)
	}
	if err != nil {
		return
	}
	code := util.RandomStr(5)
	err = o.platform.Cache.Set(code, link, time.Duration(3600)*time.Second)
	if err != nil {
		return
	}
	return &AuthLinkRsp{
		Code: code,
		Link: link,
	}, nil
}
func (o *OpenPlatform) DecryptMsg(message []byte) (*BaseMessageDecEncrypt, error) {
	mData := new(MessageEncrypt)
	if err := xml.Unmarshal(message, mData); err != nil {
		return nil, err
	}
	var data BaseMessageDecEncrypt
	_, xmlBytes, xmlErr := util.DecryptMsg(o.config.AppID, mData.Encrypt, o.config.EncodingAESKey)
	if xmlErr != nil {
		return nil, xmlErr
	}
	data.XmlByte = string(xmlBytes)
	if err := xml.Unmarshal(xmlBytes, &data); err != nil {
		return nil, err
	}
	if &data != nil {
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
		if err != nil || token == "" {
			platform.SetComponentAccessToken(verifyTicket)
		}
	}
	return platform
}

func NewOpenPlatform(isToken bool) (IOpenPlatform, error) {
	memory := cache.NewMemory()
	config := &openConfig.Config{
		AppID:          viper.GetString("wechatOpenPlatform.Appid"),
		AppSecret:      viper.GetString("wechatOpenPlatform.AppSecret"),
		Token:          viper.GetString("wechatOpenPlatform.Token"),
		EncodingAESKey: viper.GetString("wechatOpenPlatform.EncodingAESKey"),
		Cache:          memory,
	}
	platform := new(OpenPlatform)
	newWechat := wechat.NewWechat()
	platforms := newWechat.GetOpenPlatform(config)
	verifyTicket := viper.GetString("wechatOpenPlatform.VerifyTicket")
	if isToken && verifyTicket != "" {
		token, err := platforms.GetComponentAccessToken()
		fmt.Println("token get", token)
		fmt.Println("token err", err)
		if err != nil || token == "" {
			_, errs := platforms.SetComponentAccessToken(verifyTicket)
			if errs != nil {
				return nil, errs
			}
		}
	}
	platform.config = config
	platform.platform = platforms
	return platform, nil
}
