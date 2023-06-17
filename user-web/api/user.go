package api

import (
	context1 "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"time"
)

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

func GetUserList(context *gin.Context) {
	ip := "127.0.0.1"
	port := 50051

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[getUserList] 连接用户服务失败", "msg", err.Error())
	}

	//	生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	rsp, err := userSrvClient.GetUserList(context1.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
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
