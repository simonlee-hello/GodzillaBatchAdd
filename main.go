/*
批量添加webshell到Godzilla
*/
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开SQLite数据库连接
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 读取urls.txt文件
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 逐行读取urls.txt文件并插入数据库
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		id := uuid.New().String()
		createTime := time.Now().Format("2006-01-02 15:04:05")
		updateTime := createTime

		// 插入数据到main表
		_, err := db.Exec(`
			INSERT INTO shell (id, url, password, secretKey, payload, cryption, encoding, headers, reqLeft, reqRight, connTimeout, readTimeout, proxyType, proxyHost, proxyPort, remark, note, createTime, updateTime)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, url, "test", "test", "PhpDynamicPayload", "PHP_XOR_BASE64",
			"UTF-8", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\nAccept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
			"", "", 3000, 60000, "HTTP", "127.0.0.1", 8080, "note", "note", createTime, updateTime)
		if err != nil {
			log.Fatal(err)
		}

		// 插入数据到shellEnv表
		_, err = db.Exec(`
			INSERT INTO shellEnv (shellId, key, value)
			VALUES (?, ?, ?)`,
			id, "ENV_GROUP_ID", "/green")
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("数据插入完成")
}
