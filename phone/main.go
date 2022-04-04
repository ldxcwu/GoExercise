package main

import (
	"fmt"
	phonedb "phone/db"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "10.1.116.80"
	port     = 3306
	user     = "root"
	password = "123456"
	dbname   = "phone"
)

func main() {
	dsn := "root:123456@tcp(10.1.116.80:3306)/phone?charset=utf8"
	dirverName := "mysql"
	must(phonedb.Reset(dirverName, dsn, dbname))

	must(phonedb.Migrate(dirverName, dsn))

	db, err := phonedb.Open(dirverName, dsn)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Println("phone number: ", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("\tUpdating ro removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.Id))
			} else {
				p.Number = number
				must(db.UpdatePhone(p))
			}
		} else {
			fmt.Println("\tNo changes required...")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(phone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(phone, "")
}

// \d   匹配数字
// \W	匹配任意不是字母，数字，下划线，汉字的字符
// \S	匹配任意不是空白符的字符
// \D	匹配任意非数字的字符
// \B	匹配不是单词开头或结束的位置
// [^x]	匹配除了x以外的任意字符
// [^aeiou]	匹配除了aeiou这几个字母以外的任意字符
// https://deerchao.cn/tutorials/regex/regex.htm

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }
