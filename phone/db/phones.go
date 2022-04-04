package db

import (
	"database/sql"
)

//Phone represents the phone_numbers table in the DB.
type Phone struct {
	Id     int
	Number string
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.Id, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {

	data := []string{
		"1234567890 ",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil

	// _, err = insertPhone(db, "1234567890")
	// must(err)
	// _, err = insertPhone(db, "123 456 7891")
	// must(err)
	// _, err = insertPhone(db, "(123) 456 7892")
	// must(err)
	// _, err = insertPhone(db, "(123) 456-7893")
	// must(err)
	// _, err = insertPhone(db, "123-456-7894")
	// must(err)
	// id, err := insertPhone(db, "123-456-7890")
	// must(err)
	// _, err = insertPhone(db, "1234567892")
	// must(err)
	// _, err = insertPhone(db, "(123)456-7892")
	// must(err)
}

func (db *DB) UpdatePhone(p Phone) error {
	statement := "UPDATE phone_numbers set value=? where id=?"
	_, err := db.db.Exec(statement, p.Number, p.Id)
	return err
}

func (db *DB) DeletePhone(id int) error {
	statement := "DELETE FROM phone_numbers where id=?"
	_, err := db.db.Exec(statement, id)
	return err
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM phone_numbers WHERE value=?", number)
	err := row.Scan(&p.Id, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

type DB struct {
	db *sql.DB
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumberTable(db)
	if err != nil {
		return err
	}
	return db.Close()
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

func Reset(driver string, dataSource, dbName string) error {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
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
