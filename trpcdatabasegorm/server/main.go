package main

import (
	"context"
	"time"

	"trpc.group/trpc-go/trpc-database/gorm"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

// User is the model struct.
type User struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Uin        string    `gorm:"column:uin;type:varchar(128);not null;default:'';uniqueIndex:uin_unique"`
	Uid        *int64    `gorm:"column:uid;type:bigint;default:null"`
	Att        string    `gorm:"column:att;type:varchar(64);default:''"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	ModifyTime time.Time `gorm:"column:modify_time;autoUpdateTime"`
}

// TableName sets the table name for the User model.
func (User) TableName() string {
	return "t_user"
}

func main() {
	_ = trpc.NewServer()

	cli, err := gorm.NewClientProxy("trpc.mysql.server.service")
	if err != nil {
		panic(err)
	}

	// Create record
	uid := int64(1234567890)
	insertUser := User{Uin: "gorm-client", Uid: &uid, Att: "test"}
	result := cli.Create(&insertUser)
	log.Infof("inserted data's primary key: %d, err: %v", insertUser.ID, result.Error)

	// Query record
	var queryUser User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := cli.WithContext(ctx).First(&queryUser).Error; err != nil {
		log.Errorf("query user failed: %v", err)
	} else {
		log.Infof("query user: %+v", queryUser)
	}

	// Delete record
	// deleteUser := User{ID: insertUser.ID}
	// if err := cli.Delete(&deleteUser).Error; err != nil {
	// 	panic(err)
	// }
	// log.Info("delete record succeed")

	// For more use cases, see https://gorm.io/docs/create.html
}
