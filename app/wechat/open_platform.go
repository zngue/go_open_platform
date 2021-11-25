package wechat

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/openplatform"
	openConfig "github.com/silenceper/wechat/v2/openplatform/config"
	"github.com/silenceper/wechat/v2/util"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg"
	"github.com/zngue/go_open_platform/app/model"
	"io/ioutil"
	"net/http"
	"time"
)

type OpenPlatform struct {
	config   *openConfig.Config
	platform *openplatform.OpenPlatform
}

func (o *BaseMessageDecEncrypt) SaveMessage() {
	pkg.MysqlConn.Model(&model.MessageContent{}).Create(&model.MessageContent{
		Content:       o.XmlByte,
		CreateTime:    time.Now().Unix(),
		OriginContent: o.OriginXml,
	})
}

type IOpenPlatform interface {
	Platform(isToken bool) *openplatform.OpenPlatform
	DecryptMsg(message []byte) (*BaseMessageDecEncrypt, error)
	GetLinkByCode(code string) (string, error)
	AuthLink(req *AuthLinkRequest) (authLin *AuthLinkRsp, err error)
	AccountInfo(authCode string) error
	DaiLiAuth() (string, error)
	Open(ctx *gin.Context,appid string)
}

func (o *OpenPlatform) GetLinkByCode(code string) (string, error) {
	link := o.platform.Cache.Get(code)
	if linkUrl, ok := link.(string); ok {
		return linkUrl, nil
	} else {
		return "", errors.New("not exit code link")
	}
}
func (o *OpenPlatform) DaiLiAuth() (string, error) {
	officialAccount := o.platform.GetOfficialAccount("wx0372cdfcefa08b99")
	oauth := officialAccount.GetOauth()
	url := "https://api.zngue.com/authorization.php"
	return oauth.GetRedirectURL(url, "snsapi_userinfo", "STATE&")
}

func (o *OpenPlatform) Open(ctx *gin.Context,appid string){

	server := o.platform.GetOfficialAccount(appid).GetServer(ctx.Request, ctx.Writer)
	server.SetMessageHandler(func(mixMessage *message.MixMessage) *message.Reply {
		reply := &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("TESTCOMPONENT_MSG_TYPE_TEXT_callback"),
		}
		return reply
	})
	if err := server.Serve(); err != nil {
		fmt.Println(err)
		return
	}
	if err := server.Send(); err != nil {
		fmt.Println(err)
		return
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
func (o *OpenPlatform) AccountInfo(authCode string) error {
	authInfo, err := o.platform.QueryAuthCode(authCode)
	if err != nil {
		return err
	}
	info, baseInfo, errps := o.platform.GetAuthrInfo(authInfo.Appid)
	if errps != nil {
		return errps
	}
	dbCtx := pkg.MysqlConn.Model(&model.OfficialAccount{})
	var num int64
	errs := dbCtx.Where("appid = ?", baseInfo.Appid).Count(&num).Error
	if errs != nil {
		return errs
	}
	data := model.OfficialAccount{
		NickName:        info.NickName,
		Appid:           baseInfo.Appid,
		HeadImg:         info.HeadImg,
		ServiceTypeInfo: info.ServiceTypeInfo.ID,
		VerifyTypeInfo:  info.VerifyTypeInfo.ID,
		UserName:        info.UserName,
		PrincipalName:   info.PrincipalName,
		OpenStore:       info.BusinessInfo.OpenStore,
		OpenScan:        info.BusinessInfo.OpenScan,
		OpenPay:         info.BusinessInfo.OpenPay,
		OpenCard:        info.BusinessInfo.OpenCard,
		OpenShake:       info.BusinessInfo.OpenShake,
	}
	var dbErr error
	if num > 0 {
		dbErr = dbCtx.Where("appid = ?", baseInfo.Appid).Updates(&data).Error
	} else {
		dbErr = dbCtx.Create(&data).Error
	}
	return dbErr
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
	data.OriginXml = string(message)
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

type Tick struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func GetToken() (string,error) {
	response, err := http.Get("https://api.zngue.com/platform/message/ticket")
	if err != nil {
		return "", err
	}
	all, errs := ioutil.ReadAll(response.Body)
	if errs != nil {
		return "", errs
	}
	var ts Tick
	if errsok := json.Unmarshal(all, &ts); errsok != nil {
		return "", errsok
	}
	if ts.Code==200 {
		return ts.Data,nil
	}
	return "",errors.New(ts.Msg)
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
	if len(verifyTicket)==0 {
		verifyTicket, _ = GetToken()
	}
	if isToken && verifyTicket != "" {
		token, err := platforms.GetComponentAccessToken()
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
