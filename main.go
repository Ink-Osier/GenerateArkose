package main

import (
    "github.com/acheong08/funcaptcha"
    "github.com/gin-gonic/gin"
)

type TokenRequest struct {
    Type string `json:"type"`
}

func main() {
    router := gin.Default()

    router.POST("/api/arkose/token", func(c *gin.Context) {
        var req TokenRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "bad request"})
            return
        }

        // 创建Solver实例
        solver := funcaptcha.NewSolver()

        // 加载HAR文件，这里需要提供包含HAR文件的目录路径
        funcaptcha.WithHarpool(solver)

        // 根据请求的类型选择arkType
        arkType := funcaptcha.ArkVerAuth
        switch req.Type {
        case "gpt-4":
            arkType = funcaptcha.ArkVerChat4
        case "gpt-3":
            arkType = funcaptcha.ArkVerChat3
        default:
            arkType = funcaptcha.ArkVerAuth
        }

        // 调用GetOpenAIToken方法
        token, err := solver.GetOpenAIToken(arkType, "")
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        // 将token放到JSON响应的token字段里
        c.JSON(200, gin.H{"token": token})
    })

    router.Run(":8080")
}
