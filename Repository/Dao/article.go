package Dao

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`
	//作者ID
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}
