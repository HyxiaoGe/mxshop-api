package api

import (
	context1 "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	// 将 grpc 的code 转换成http 状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if ok {
				switch e.Code() {
				case codes.NotFound:
					ctx.JSON(http.StatusNotFound, gin.H{
						"msg": e.Message(),
					})
				case codes.Internal:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": e.Message(),
					})
				case codes.InvalidArgument:
					ctx.JSON(http.StatusBadRequest, gin.H{
						"msg": e.Message(),
					})
				case codes.Unavailable:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "用户服务不可用",
					})
				default:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "其他错误",
					})
				}
				return
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(context *gin.Context) {

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[getUserList] 连接用户服务失败", "msg", err.Error())
	}

	//	生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	pn := context.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := context.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context1.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[getUserList] 获取用户列表失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, context)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	context.JSON(http.StatusOK, result)

}

func PassWordLogin(ctx *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": removeTopStruct(errs.Translate(global.Trans)),
		})
		return
	}

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接用户服务失败", "msg", err.Error())
	}

	//	生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	// 登录的逻辑
	if rsp, err := userSrvClient.GetUserByMobile(context1.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			}
			return
		}
	} else {
		// 只是查询到了用户信息，还需要比对密码
		if pasRsp, pasErr := userSrvClient.CheckPassWord(context1.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); pasErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败"})
		} else {
			if pasRsp.Success {
				ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "登录失败"})
			}
		}
	}

}
