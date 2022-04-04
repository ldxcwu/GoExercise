package main

import (
	"database/sql"
	"fmt"
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
	// dsn := "root:123456@tcp(10.1.116.80:3306)/phone?charset=utf8"
	// db, err := sql.Open("mysql", dsn)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname))
	must(err)

	// err = resetDB(db, dbname)
	// must(err)

	defer db.Close()
	// must(db.Ping())

	must(createPhoneNumberTable(db))

	_, err = insertPhone(db, "1234567890")
	must(err)
	_, err = insertPhone(db, "123 456 7891")
	must(err)
	_, err = insertPhone(db, "(123) 456 7892")
	must(err)
	_, err = insertPhone(db, "(123) 456-7893")
	must(err)
	_, err = insertPhone(db, "123-456-7894")
	must(err)
	id, err := insertPhone(db, "123-456-7890")
	must(err)
	_, err = insertPhone(db, "1234567892")
	must(err)
	_, err = insertPhone(db, "(123)456-7892")
	must(err)

	number, _ := getPhone(db, id)
	fmt.Println("phone number: ", number)

	phones, err := allPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Println("phone number: ", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("\tUpdating ro removing...", number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				must(deletePhone(db, p.id))
			} else {
				p.number = number
				must(updatePhone(db, p))
			}
		} else {
			fmt.Println("\tNo changes required...")
		}
	}
}

func updatePhone(db *sql.DB, p phone) error {
	statement := "UPDATE phone_numbers set value=? where id=?"
	_, err := db.Exec(statement, p.number, p.id)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := "DELETE FROM phone_numbers where id=?"
	_, err := db.Exec(statement, id)
	return err
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE value=?", number)
	err := row.Scan(&p.id, &p.number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	//insert data VALUES(could contains more than 1 ? )
	//just make sure that  relating args behind
	statement := `INSERT INTO phone_numbers(value) VALUES(?)`
	ret, err := db.Exec(statement, phone)
	if err != nil {
		return -1, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

type phone struct {
	id     int
	number string
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id=?", id)
	// row := db.QueryRow("SELECT value FROM phone_numbers WHERE id=?", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}

	return number, nil
}

func createPhoneNumberTable(db *sql.DB) error {

	tryDropPhoneNumberTable(db)

	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id INT(11) AUTO_INCREMENT PRIMARY KEY,
			value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	return err
}

func tryDropPhoneNumberTable(db *sql.DB) error {
	statement := "DROP TABLE phone_numbers"
	_, err := db.Exec(statement)
	return err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	must(err)
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	must(err)
	return nil
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
