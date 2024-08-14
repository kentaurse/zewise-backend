/*
Package main - ZeWise 后端服务入口
该文件用于定义 WeWise 后端服务入口
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"zewise.space/backend/configs"
	"zewise.space/backend/consts"
	"zewise.space/backend/controllers"
	"zewise.space/backend/middlewares"
	"zewise.space/backend/models"
	"zewise.space/backend/services"
	"zewise.space/backend/stores"
	"zewise.space/backend/utils/functools"
)

var (
	config            *configs.Config
	redisClient       *redis.Client
	mongoClient       *mongo.Client
	minioClient       *minio.Client
	storage           *stores.Storage
	controllerFactory *controllers.Factory
	middlewareFactory *middlewares.Factory
)

func init() {
	var err error

	// 初始化配置文件
	config, err = configs.NewConfig()
	if err != nil {
		panic(err)
	}

	// 初始化 Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}

	// 初始化 MongoDB
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(
		functools.JoinStrings("mongodb://", config.MongoDB.Host, ":", fmt.Sprint(config.MongoDB.Port)),
	))
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	// 初始化 MinIO
	minioClient, err = minio.New(
		functools.JoinStrings(config.MinIO.Host, ":", fmt.Sprint(config.MinIO.Port)),
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.MinIO.Access, config.MinIO.Secret, ""),
			Secure: false,
		},
	)
	if err != nil {
		panic(err)
	}
	err = models.SetupBucket(minioClient)
	if err != nil {
		panic(err)
	}

	// 初始化存储
	storage = stores.NewStore(redisClient, mongoClient, config.MongoDB.DBName, minioClient)

	// 初始化控制器工厂
	controllerFactory = controllers.NewFactory(
		services.NewService(storage),
	)
	// 初始化中间件工厂
	middlewareFactory = middlewares.NewFactory(storage)
}

func main() {
	// 配置 Fiber
	var fiberConfig fiber.Config
	switch config.Env.Type {
	case "development":
		fiberConfig.Prefork = false
	case "production":
		fiberConfig.Prefork = true
	}
	fiberConfig.BodyLimit = consts.REQUEST_BODY_LIMIT

	// 创建 Fiber 实例
	app := fiber.New(fiberConfig)

	// 设置 Fiber 中间件
	app.Use(logger.New(logger.Config{
		Format: "[${time}][${latency}][${status}][${method}] ${path}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: config.Compress.Level,
	}))

	// Auth 中间件
	auth := middlewareFactory.NewTokenAuthMiddleware()

	// api 路由
	api := app.Group("/api")

	// Auth 路由
	authController := controllerFactory.NewAuthController()
	authGroup := api.Group("/auth")
	authGroup.Post("/login", authController.NewLoginHandler())                                // 登录
	authGroup.Post("/logout", auth.NewMiddleware(), authController.NewLogoutHandler())        // 登出
	authGroup.Post("/refresh", auth.NewMiddleware(), authController.NewRefreshTokenHandler()) // 刷新令牌
	// authGroup.Post("/verify/mail")

	// User 路由
	userController := controllerFactory.NewUserController()
	user := api.Group("/user")
	user.Get("/profile", userController.NewProfileHandler())                                       // 获取用户资料信息
	user.Post("/register", userController.NewRegisterHandler())                                    // 注册
	user.Post("/update/profile", auth.NewMiddleware(), userController.NewUpdateProfileHandler())   // 更新用户资料
	user.Post("/update/avatar", auth.NewMiddleware(), userController.NewUpdateAvatarHandler())     // 更新用户头像
	user.Post("/update/password", auth.NewMiddleware(), userController.NewUpdatePasswordHandler()) // 更新用户密码

	panic(app.Listen(functools.JoinStrings(config.Server.Host, ":", fmt.Sprint(config.Server.Port))))
}
