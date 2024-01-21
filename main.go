package main

import (
    "fmt"
    "context"
    "github.com/acheong08/funcaptcha"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "net/http"
    "os"
    "strconv"
    "time"
    "strings"
)

var redisClient *redis.Client
var ctx = context.Background()

func initRedis() {
    // 从环境变量获取 Redis 的地址和端口
    redisHost := os.Getenv("REDIS_HOST")
    if redisHost == "" {
        redisHost = "redis" // 如果没有设置，默认为 localhost
    }

    redisPort := os.Getenv("REDIS_PORT")
    if redisPort == "" {
        redisPort = "6379" // 如果没有设置，默认为 6379 端口
    }

    redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

    redisDBStr := os.Getenv("REDIS_DB")
    redisDB := 0 // 默认为 0 号数据库
    if redisDBStr != "" {
        var err error
        redisDB, err = strconv.Atoi(redisDBStr)
        if err != nil {
            fmt.Println("Error parsing REDIS_DB:", err)
        }
    }

    redisPassword := os.Getenv("REDIS_PASSWORD")

    // 创建 Redis 客户端的配置
    opts := &redis.Options{
        Addr: redisAddr,
        DB:   redisDB,
    }

    if redisPassword != "" {
        opts.Password = redisPassword
    }

    // 创建一个新的 Redis 客户端
    redisClient = redis.NewClient(opts)
}

func checkAndGenerateTokens() {
    for {
        // 获取环境变量中设置的 token 阈值
        minTokenCount, _ := strconv.Atoi(os.Getenv("MIN_TOKEN_COUNT"))

        // 从环境变量获取 token 过期时间
        tokenExpiryMinutes, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRY_MINUTES"))
        if tokenExpiryMinutes == 0 {
            tokenExpiryMinutes = 25  // 如果没有设置或设置错误，默认为 25 分钟
        }

        // 获取 Redis 中的 token 数量
        size, err := redisClient.DBSize(ctx).Result()
        if err != nil {
            fmt.Println("Error getting token count from Redis:", err)
            continue
        }

        // 如果 token 数量小于阈值，则生成并存储新的 token
        if int(size) < minTokenCount {
            message := fmt.Sprintf("Current token count: %d is less than min token count: %d, generating new tokens...", size, minTokenCount)
            fmt.Println(message)

            solver := funcaptcha.NewSolver()
            funcaptcha.WithHarpool(solver)
            token, err := solver.GetOpenAIToken(funcaptcha.ArkVerChat4, "")
            if err != nil {
                fmt.Println("Error generating token:", err)
                continue
            }

            // 将 token 存入 Redis，并设置过期时间
            err = redisClient.Set(ctx, "token:"+token, token, time.Duration(tokenExpiryMinutes)*time.Minute).Err()
            if err != nil {
                fmt.Println("Error saving token to Redis:", err)
            }
        }

        // 每秒检查一次
        time.Sleep(1 * time.Second)
    }
}

func main() {
    router := gin.Default()
    initRedis()

    // 启动后台 goroutine 来检查和生成 token
    go checkAndGenerateTokens()

    router.POST("/api/arkose/token", func(c *gin.Context) {
        // 尝试从 Redis 获取一个 token
        key, err := redisClient.RandomKey(ctx).Result()
        if err != nil || key == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "no token available"})
            return
        }
    
        // 从 Redis 中删除获取到的 token
        _, delErr := redisClient.Del(ctx, key).Result()
        if delErr != nil {
            fmt.Println("Error deleting token from Redis:", delErr)
        }
    
        // 去除 "token:" 前缀
        token := strings.TrimPrefix(key, "token:")
    
        // 返回获取到的 token
        c.JSON(http.StatusOK, gin.H{"token": token})
    })
    
    

    router.Run(":8080")
}
