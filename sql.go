package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB //是一个连接池对象 因为要在其他insert等函数内使用db所以要声明为全局变量

type depot struct {
	depot_id string
	volume   int
}

type cloth struct {
	cloth_id string
	size     string
	price    int
	kinds    string
}

type produce struct {
	cloth_id string
	mer_id   string
	grade    string
}

type merchant struct {
	mer_id string
	name   string
}

func initDB() (err error) {
	// DSN:Data Source Name
	//
	dsn := "root:56874312zj@tcp(127.0.0.1:3306)/hsp"
	//dsn := "root:56874312zj@tcp(10.17.108.145:3306)/hsp"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("dsn failed")
		//panic(err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("open failed")
		return
	}

	return
}

func queryone(depot_id string) {
	var u1 depot

	//sqlStr := "select depot_id,volume from depot;"
	//之前一直查出来数据是0的原因---没有where这句使得(sqlStr, 1)里面的1无效，这里的1就是去替换?
	sqlStr := `select depot_id,volume from depot where volume > ?;`
	//sqlStr := `select depot_id,volume from depot ;`
	rowObj := db.QueryRow(sqlStr, 1)
	rowObj.Scan(&u1.depot_id, &u1.volume)
	fmt.Printf("u1:%#v\n", u1)

}

func queryMore(n int) {

	sqlStr := `select depot_id,volume from depot where volume > ?;`
	//sqlStr1 := `select size,price from cloth where size = s AND price < 100;`
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Printf("exec %s failed %v", sqlStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u1 depot
		err := rows.Scan(&u1.depot_id, &u1.volume)
		if err != nil {
			fmt.Printf("scan failed %v", sqlStr, err)
		}
		fmt.Printf("u1:%#v\n", u1)
	}

}

func queryCloth() {
	var c2 cloth

	sqlStr1 := `select * from cloth where size = 'S' AND price < 100;`
	rows, err := db.Query(sqlStr1)
	if err != nil {
		fmt.Printf("exec %s failed %v", sqlStr1, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c2.cloth_id, &c2.size, &c2.price, &c2.kinds)
		if err != nil {
			fmt.Printf("scan %s failed %v", sqlStr1, err)
		}
		fmt.Printf("c2:%#v\n", c2)
	}

}

func queryDepot() {
	var d3 depot

	sqlStr := `select * from depot order by volume DESC;`
	row := db.QueryRow(sqlStr)
	err := row.Scan(&d3.depot_id, &d3.volume)

	if err != nil {
		fmt.Printf("scan %s failed %v", sqlStr, err)
	}
	fmt.Printf("d3:%#v\n", d3)

}

func queryAcloth() {
	var c4 cloth

	sqlStr := `select * from cloth where cloth_id like 'A%';`
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("exec %s failed %v", sqlStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c4.cloth_id, &c4.size, &c4.price, &c4.kinds)
		if err != nil {
			fmt.Printf("scan %s failed %v", sqlStr, err)
		}
		fmt.Printf("c4:%#v\n", c4)
	}

}

func queryGradeMer() {
	var m5 merchant
	//注意进行联结，并指定mer_id是哪个表的 操作方法是加 表名.
	sqlStr := `select merchant.mer_id,name FROM merchant LEFT JOIN produce 
	ON merchant.mer_id = produce.mer_id where grade > 'D';`
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("exec %s failed %v", sqlStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&m5.mer_id, &m5.name)
		if err != nil {
			fmt.Printf("scan %s failed %v", sqlStr, err)
		}
		fmt.Printf("m5:%#v\n", m5)
	}

}

func insert() {
	sqlStr := `insert into depot(depot_id,volume) values("789",300);`
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("insert failed %v", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRow() {
	var c6 cloth
	sqlStr1 := `select price from cloth where size = 'S' ;`
	rows, err := db.Query(sqlStr1)
	if err != nil {
		fmt.Printf("exec %s failed %v", sqlStr1, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c6.price)
		if err != nil {
			fmt.Printf("scan %s failed %v", sqlStr1, err)
		}
		var newPrice float32
		newPrice = float32(c6.price) * 1.10
		sqlStr := "update cloth set price = ? where size = 'S'"

		ret, err := db.Exec(sqlStr, newPrice)
		if err != nil {
			fmt.Printf("update failed, err:%v\n", err)
			return
		}
		n, err := ret.RowsAffected() // 操作影响的行数
		if err != nil {
			fmt.Printf("get RowsAffected failed, err:%v\n", err)
			return
		}
		fmt.Printf("update success, affected rows:%d\n", n)
	}

}

// 删除数据
func deleteRow() {
	sqlStr := "delete from produce where grade > 'D'"
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init failed")
	}
	fmt.Println("连接数据库成功")

	//queryone("123")
	//queryMore(0)

	//queryCloth()
	//queryDepot()
	//queryAcloth()
	//queryGradeMer()
	//updateRow()
	deleteRow()
	//inert()
	//updateRow()
	//deleteRow()
}
