package blog

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	. "mymanager/models/common"
)

type blogSyncAoModel struct {
	Model
	BlogSyncDb     BlogSyncDbModel
	BlogSyncAutoDb BlogSyncAutoDbModel
	BlogCsdnAo     BlogCsdnAoModel
	BlogGitAo      BlogGitAoModel
}

func (this *blogSyncAoModel) SearchAuto(userId int, where BlogSyncAuto, limit CommonPage) BlogSyncAutos {
	where.UserId = userId
	return this.BlogSyncAutoDb.Search(where, limit)
}

func (this *blogSyncAoModel) GetAuto(userId int, blogSyncAutoId int) BlogSyncAuto {
	cardInfo := this.BlogSyncAutoDb.Get(blogSyncAutoId)
	if cardInfo.UserId != userId {
		Throw(1, "你没有该权限")
	}
	return cardInfo
}

func (this *blogSyncAoModel) DelAuto(userId int, blogSyncAutoId int) {
	this.GetAuto(userId, blogSyncAutoId)

	this.BlogSyncAutoDb.Del(blogSyncAutoId)
}

func (this *blogSyncAoModel) AddAuto(userId int, blogSyncAuto BlogSyncAuto) {
	blogSyncAuto.UserId = userId
	this.BlogSyncAutoDb.Add(blogSyncAuto)
}

func (this *blogSyncAoModel) ModAuto(userId int, blogSyncAutoId int, blogSyncAuto BlogSyncAuto) {
	this.GetAuto(userId, blogSyncAutoId)

	blogSyncAuto.UserId = userId
	this.BlogSyncAutoDb.Mod(blogSyncAutoId, blogSyncAuto)
}

func (this *blogSyncAoModel) SearchTask(userId int, where BlogSync, limit CommonPage) BlogSyncs {
	where.UserId = userId
	return this.BlogSyncDb.Search(where, limit)
}

func (this *blogSyncAoModel) AddTask(userId int, accessToken string, gitUrl string, syncType int) {
	data := BlogSync{
		UserId:       userId,
		AccessToken:  accessToken,
		GitUrl:       gitUrl,
		SyncType:     syncType,
		State:        BlogStateEnum.STATE_BEGIN,
		StateMessage: "",
	}
	syncId := this.BlogSyncDb.Add(data)
	this.Queue.Produce("blog_sync", syncId)
}

func (this *blogSyncAoModel) GetTask(userId int, blogSyncId int) BlogSync {
	data := this.BlogSyncDb.Get(blogSyncId)
	if data.UserId != userId {
		Throw(1, "权限不足")
	}
	return data
}

func (this *blogSyncAoModel) RestartTask(userId int, blogSyncId int) {
	data := this.GetTask(userId, blogSyncId)
	if data.State != BlogStateEnum.STATE_FAIL {
		Throw(1, "非失败任务不能重启")
	}
	this.modState(blogSyncId, BlogStateEnum.STATE_BEGIN, "")

	this.Queue.Produce("blog_sync", blogSyncId)
}

func (this *blogSyncAoModel) modState(blogSyncId int, state int, stateMessage string) {
	modData := BlogSync{
		State:        state,
		StateMessage: stateMessage,
	}
	this.BlogSyncDb.Mod(blogSyncId, modData)
}

func (this *blogSyncAoModel) sync(blogSyncId int) {
	defer CatchCrash(func(e Exception) {
		this.modState(blogSyncId, BlogStateEnum.STATE_FAIL, e.GetMessage())
		panic(e.Error())
	})
	//获取同步信息
	data := this.BlogSyncDb.Get(blogSyncId)

	//执行同步
	updateProgress := func(message string) {
		this.modState(blogSyncId, BlogStateEnum.STATE_PROGRESS, message)
	}
	blogs := this.BlogGitAo.Get(data.GitUrl, updateProgress)
	this.BlogCsdnAo.Sync(data.AccessToken, data.SyncType, blogs, updateProgress)

	this.modState(blogSyncId, BlogStateEnum.STATE_SUCCESS, "成功")
}

func (this *blogSyncAoModel) syncAuto() {
	data := this.BlogSyncAutoDb.GetAll()
	for _, singleData := range data {
		this.AddTask(
			singleData.UserId,
			singleData.AccessToken,
			singleData.GitUrl,
			BlogSyncTypeEnum.TYPE_INCREMENTAL_UPDATE,
		)
	}
}

func init() {
	InitDaemon(func(this *blogSyncAoModel) {
		this.Queue.Consume("blog_sync", this.sync)
		this.Timer.Cron("0 0 23 * *", this.syncAuto)
	})
}
