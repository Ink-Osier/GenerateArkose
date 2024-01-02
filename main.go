package main

import (
    "fmt"
    "github.com/acheong08/funcaptcha"
    "github.com/gin-gonic/gin"
    "net/http"
)

type TokenRequest struct {
    Type string `form:"type"`
}

func main() {
    router := gin.Default()

    router.POST("/api/arkose/token", func(c *gin.Context) {
        var req TokenRequest
        // 使用 ShouldBind 方法来绑定 x-www-form-urlencoded 数据
        if err := c.ShouldBind(&req); err != nil {
            fmt.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
            return
        }

        // 创建 Solver 实例
        solver := funcaptcha.NewSolver()

        // 加载 HAR 文件，这里需要提供包含 HAR 文件的目录路径
        funcaptcha.WithHarpool(solver)

        // 根据请求的类型选择 arkType
        arkType := funcaptcha.ArkVerAuth
        switch req.Type {
        case "gpt-4":
            arkType = funcaptcha.ArkVerChat4
        case "gpt-3":
            arkType = funcaptcha.ArkVerChat3
        default:
            arkType = funcaptcha.ArkVerAuth
        }

        // 调用 GetOpenAIToken 方法
        token, err := solver.GetOpenAIToken(arkType, "")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // 将 token 放到 JSON 响应的 token 字段里
        c.JSON(http.StatusOK, gin.H{"token": token})
    })

    router.Run(":8080")
}
