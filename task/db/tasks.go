package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// return db.Update(func(tx *bolt.Tx) error {
	// 	_, err := tx.CreateBucketIfNotExists(taskBucket)
	// 	return err
	// })
	//Update()开启一个读写事务
	//Update()接收一个函数，表明该事务要做的事情，
	//这里首先要做的事情就是创建一个Bucket
	//Bucket represents a collection of key/value pairs inside the databse.
	return db.Update(createBucket)
}

//也可以take in an task struct，但是那样由用户控制task.key，不好控制
func CreateTask(task string) (int, error) {
	//创建任务，投入Bucket
	//1. 开启一个事务去做此事
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		//1. 获取之前创建的Bucket
		b := tx.Bucket(taskBucket)
		//NextSequence returns an autoincrementing integer for the bucket.
		//这里不处理错误，因为出错的原因只能是事务的关闭等与事务相关联的错误
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		//b.Put接收两个byte slice
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	//db.View starts a read-only transaction
	err := db.View(func(tx *bolt.Tx) error {
		//1. Get Bucket
		b := tx.Bucket(taskBucket)
		//2. Get Cursor which can be seen as a point in the bucket
		c := b.Cursor()
		//3. iterate the bucket element
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func createBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(taskBucket)
	return err
}

//itob 将 int 转为 byte 切片
// 1    -> [0, 0, 0, 0, 0, 0, 0, 1]
// 255  -> [0, 0, 0, 0, 0, 0, 0, 255]
// 256  -> [0, 0, 0, 0, 0, 0, 1, 0]
// 257  -> [0, 0, 0, 0, 0, 0, 1, 1]
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
