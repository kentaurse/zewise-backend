/*
Package stores - ZeWise 后端服务器数据访问层
该文件用于声明存储对象类
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package stores

import (
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage 存储对象
type Storage struct {
	redis       *redis.Client // redis 客户端
	mongo       *mongo.Client // mongo 客户端
	minio       *minio.Client // minio 客户端
	AuthStorage *AuthStorage  // 认证相关存储
	UserStorage *UserStorage  // 用户相关存储
	// PostStore    *PostStore    // 文章相关存储
	// CommentStore *CommentStore // 评论相关存储
	// ReplyStore   *ReplyStore   // 回复相关存储
}

/*
NewStore 新建存储

参数：
  - redis：redis 客户端
  - mongo：mongo 客户端
  - mongoDBName：mongo 数据库名称

返回：
  - *Storage：存储对象
*/
func NewStore(redis *redis.Client, mongo *mongo.Client, mongoDBName string, minio *minio.Client) *Storage {
	mongoDataBase := mongo.Database(mongoDBName)
	return &Storage{
		redis:       redis,
		mongo:       mongo,
		minio:       minio,
		AuthStorage: &AuthStorage{redis, mongoDataBase},
		UserStorage: &UserStorage{redis, mongoDataBase, minio},
		// PostStore:    &PostStore{redis, mongoDataBase},
		// CommentStore: &CommentStore{redis, mongoDataBase},
		// ReplyStore:   &ReplyStore{redis, mongoDataBase},
	}
}

/*
NewSession 新建事务

返回：
  - mongo.Session：mongo 事务对象
  - error：错误信息
*/
func (storage *Storage) NewSession() (mongo.Session, error) {
	return storage.mongo.StartSession()
}
