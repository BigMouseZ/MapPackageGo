package mapinit

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var Config map[string]string

func MapInit() {
	//初始化配置
	Config = InitConfig("./mappackagego.properties")
	if Config == nil {
		return
	}
	Dbinit()
}

//读取key=value类型的配置文件
func InitConfig(path string) map[string]string {
	config := make(map[string]string)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println("打包配置不存在 :" + path)
		return nil
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("打包配置异常 :", err)
			return nil
		}
		s := strings.TrimSpace(string(b))
		indexSkip := strings.Index(s, "#")
		if indexSkip >= 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}

var (
	Db *sql.DB
)

func Dbinit() *sql.DB {
	dbUrl := Config["datasource_url"]
	var err error
	DbInit,err := sql.Open("mysql",dbUrl)
	if err != nil{
		panic("连接数据库失败:" + err.Error())
	}else{
		fmt.Println("connect to mysql success")
	}

	Db = DbInit
	return DbInit
}
