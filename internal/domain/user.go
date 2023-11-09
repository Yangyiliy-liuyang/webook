package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	//UTC 0 的时区
	Ctime time.Time
	//Addr     Address
	Nickname string
	Birthday time.Time
	AboutMe  string
}

//type Address struct {
//	Province string
//	Region   string
//}

/*func (u User) ValidateEmail() bool {
	//正则表达式
	return u.Email
}*/
