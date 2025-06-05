package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hobbyGG/RestfulAPI_forum/contants/code"
	"github.com/hobbyGG/RestfulAPI_forum/contants/contant"
	"github.com/hobbyGG/RestfulAPI_forum/contants/errors"
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"github.com/hobbyGG/RestfulAPI_forum/service"
	"go.uber.org/zap"
)

func CreatePostHandler(ctx *gin.Context) {
	// 用户创建帖子
	// 处理数据
	postParam := new(models.ParamCreatePost)
	if err := ctx.ShouldBindJSON(&postParam); err != nil {
		zap.L().Error("ShouldBindJSON error", zap.Error(err))
		ResponseError(ctx, code.InvalidParam)
		return
	}
	uid := ctx.GetInt64("uid")
	if postParam.CommID == 0 {
		postParam.CommID = 1
	}

	// service层实现创建帖子
	var postID int64
	var err error
	if postID, err = service.CreatePost(postParam, uid); err != nil {
		zap.L().Error("service.CreatePost", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}

	ResponseSuccess(ctx, postID)
}

func GetPostHander(ctx *gin.Context) {
	// 通过url获取post的id参数
	postID := ctx.Param(contant.URLParamID)
	if postID == "" {
		zap.L().Error("ctx.Param error", zap.Error(errors.ErrAuthType))
		ResponseError(ctx, code.InvalidParam)
		return
	}
	var (
		err  error
		post *models.Post
	)
	if post, err = service.GetPost(postID); err != nil {
		zap.L().Error("service.GetPost error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}
	ResponseSuccess(ctx, post)
}

func GetPostsHander(ctx *gin.Context) {
	// 解析query参数
	pageStr, ok := ctx.GetQuery(contant.KeyPageStr)
	if !ok {
		pageStr = "0"
	}
	sizeStr, ok := ctx.GetQuery(contant.KeySizeStr)
	if !ok {
		sizeStr = "10"
	}
	sorted, ok := ctx.GetQuery(contant.KeySortedStr)
	if !ok {
		sorted = "time"
	}
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	posts, err := service.GetPosts(page, size, sorted)
	if err != nil {
		zap.L().Error("service.GetPosts error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}

	ResponseSuccess(ctx, posts)
}

func PostVoteHandler(ctx *gin.Context) {
	// 处理参数
	voteParam := new(models.ParamVote)
	if err := ctx.ShouldBindJSON(voteParam); err != nil {
		zap.L().Error("ctx.ShouldBindJSON error", zap.Error(err))
		ResponseError(ctx, code.InvalidParam)
		return
	}
	pid := voteParam.PostID
	vote := voteParam.Vote
	uid := ctx.GetInt64(contant.StrUID)

	// 进入服务
	if err := service.PostVote(pid, uid, vote); err != nil {
		zap.L().Error("service.PostVote error", zap.Error(err))
		ResponseError(ctx, code.ServeBusy)
		return
	}

	ResponseSuccess(ctx, nil)
}
