package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if duplicateErr == me.Number {
			// todo 用户冲突，邮箱 唯一索引冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

type User struct {
	// TODO 为什么使用自增主键？
	//数据库中的数据存储是一个树型结构，自增意味着树朝一个方向增长，id相邻的大概率在磁盘上也是相邻的
	//，充分利用操作系统预读机制。
	//不是自增则意味中间插入数据，页分页
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	//TODO 为什么不用time.time : UTC 0 的时区
	// 整个系统内部都使用UTC 0 的时区，
	// 在要返回给前端的时候才改成UTF8 或者直接交给前端处理
	//服务器 go应用 数据库

	Ctime int64
	Utime int64

	//json
	//Addr string
}

//type Address struct {
//	uid int
//}
