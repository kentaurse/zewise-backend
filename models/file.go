/*
Package models - ZeWise 数据库模型
该文件用于声明对象存储模型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"context"

	"github.com/minio/minio-go/v7"
)

const USER_AVATAR_BUCKET = "avatars"

/*
SetupBucket 初始化存储桶

参数：
  - *client：MinIO 客户端

返回：
  - error：错误信息
*/
func SetupBucket(client *minio.Client) error {
	err := client.MakeBucket(context.TODO(), USER_AVATAR_BUCKET, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), USER_AVATAR_BUCKET)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}
