package controllers

import (
	. "mymanager/models/common"
	. "mymanager/models/category"
	. "mymanager/models/user"
)

type CategoryController struct {
	BaseController
	CategoryAo CategoryAoModel
	UserLoginAo UserLoginAoModel
}

func (this *CategoryController) Search_Json() interface{} {
	//检查输入参数
	var where Category
	this.CheckGet(&where)

	var limit CommonPage
	this.CheckGet(&limit)

	//检查权限
	loginUser := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	return this.CategoryAo.Search(loginUser.UserId,where, limit)
}

func (this *CategoryController) Get_Json() interface{} {
	//检查输入参数
	var category Category
	this.CheckGet(&category)

	//检查权限
	loginUser := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	return this.CategoryAo.Get(loginUser.UserId,category.CategoryId)
}

func (this *CategoryController) Add_Json() {
	//检查输入参数
	var category Category
	this.CheckPost(&category)

	//检查权限
	loginUser := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	this.CategoryAo.Add(loginUser.UserId,category)
}

func (this *CategoryController) Del_Json() {
	//检查输入参数
	var category Category
	this.CheckPost(&category)

	//检查权限
	loginUser := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	this.CategoryAo.Del(loginUser.UserId,category.CategoryId)
}

func (this *CategoryController) Mod_Json() {
	//检查输入参数
	var category Category
	this.CheckPost(&category)

	//检查权限
	loginUser := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	this.CategoryAo.Mod(loginUser.UserId, category.CategoryId,category)
}
