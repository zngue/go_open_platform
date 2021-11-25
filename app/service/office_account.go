package service

import (
	"github.com/zngue/go_helper/pkg"
	"github.com/zngue/go_open_platform/app/model"
)

type IOfficialAccount interface {
	Add(request *OfficialAccountRequest) error
	Save(request *OfficialAccountRequest) error
	List(request *OfficialAccountRequest) (*[]model.OfficialAccount, error)
	Detail(request *OfficialAccountRequest) (*model.OfficialAccount, error)
	Delete(request *OfficialAccountRequest) error
	GetModel() interface{}
}
type OfficialAccount struct {
}
type OfficialAccountRequest struct {
	pkg.CommonRequest
	ID int `form:"id" field:"id" where:"eq" default:"0"`
	Appid string `form:"appid" field:"appid" where:"eq" default:""`
}

func NewOfficialAccount() IOfficialAccount {
	return new(OfficialAccount)
}
func (i *OfficialAccount) GetModel() interface{} {
	return new(model.OfficialAccount)
}

// Add 添加
func (i *OfficialAccount) Add(request *OfficialAccountRequest) error {
	request.ReturnType = 3
	return pkg.MysqlConn.Model(i.GetModel()).Create(request.Data).Error
}

// Save 修改
func (i *OfficialAccount) Save(request *OfficialAccountRequest) error {
	request.ReturnType = 3
	return request.Init(pkg.MysqlConn.Model(i.GetModel()), *request).Updates(request.Data).Error
}
func (i *OfficialAccount) List(request *OfficialAccountRequest) (*[]model.OfficialAccount, error) {
	var list []model.OfficialAccount
	err := request.Init(pkg.MysqlConn.Model(i.GetModel()), *request).Find(&list).Error
	return &list, err
}

// Detail 详情
func (i *OfficialAccount) Detail(request *OfficialAccountRequest) (*model.OfficialAccount, error) {
	var detail model.OfficialAccount
	request.ReturnType = 3
	err := request.Init(pkg.MysqlConn.Model(i.GetModel()), *request).First(&detail).Error
	return &detail, err
}

// Delete 删除
func (i *OfficialAccount) Delete(request *OfficialAccountRequest) error {
	request.ReturnType = 3
	return request.Init(pkg.MysqlConn, *request).Delete(i.GetModel()).Error
}