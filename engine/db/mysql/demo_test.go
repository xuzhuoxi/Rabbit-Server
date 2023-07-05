// Package mysql
// Create on 2023/7/3
// @author xuzhuoxi
package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// TT
// 定义一个结构体来存储TT表的数据
type TT struct {
	id   int
	name string
	age  int
}

func main() {
	// 创建一个数据库连接对象
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	defer db.Close() // 延迟关闭数据库连接

	// 从数据库中读取表名为TT的全部数据
	rows, err := db.Query("select id, name, age from TT")
	if err != nil {
		fmt.Println("query failed,", err)
		return
	}
	defer rows.Close() // 延迟关闭结果集

	// 定义一个切片来存储查询结果
	var tts []TT

	// 遍历结果集，将每一行数据存入切片中
	for rows.Next() {
		var tt TT
		err := rows.Scan(&tt.id, &tt.name, &tt.age)
		if err != nil {
			fmt.Println("scan failed,", err)
			return
		}
		tts = append(tts, tt)
	}

	// 打印查询结果
	fmt.Println("query result:")
	for _, tt := range tts {
		fmt.Printf("id: %d, name: %s, age: %d\n", tt.id, tt.name, tt.age)
	}

	// 修改id=1001的数据，并写回到数据库
	stmt, err := db.Prepare("update TT set name = ?, age = ? where id = ?")
	if err != nil {
		fmt.Println("prepare failed,", err)
		return
	}
	defer stmt.Close() // 延迟关闭预编译语句

	res, err := stmt.Exec("Alice2", 25, 1001) // 执行预编译语句，传入参数
	if err != nil {
		fmt.Println("exec failed,", err)
		return
	}

	row, err := res.RowsAffected() // 获取影响的行数
	if err != nil {
		fmt.Println("rows affected failed,", err)
		return
	}

	fmt.Println("update success:", row) // 打印更新成功的信息

}
