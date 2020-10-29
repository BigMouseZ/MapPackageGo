package main

import (
	"fmt"
	"github.com/beevik/etree"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"time"

	"MapPackageGo/mapinit"
)

func main() {
	start := time.Now() // 获取当前时间
	mapinit.MapInit()
	pak_name := mapinit.Config["pak_name"]
	map_type := mapinit.Config["map_type"]
	zoom_min := mapinit.Config["zoom_min"]
	zoom_max := mapinit.Config["zoom_max"]
	gmapnetcache := mapinit.Config["tablename"]
	sql := " select Type, Zoom, X, Y, Tile from " + gmapnetcache + "  where Type = " + map_type + " and Zoom >= " + zoom_min + " and Zoom <= " + zoom_max + " ORDER BY zoom,x,y"
	rows, err := mapinit.Db.Query(sql)
	if err != nil {
		fmt.Printf("select fail [%s]", err)
	}
	exist, err := PathExists("./" + pak_name + ".pak")
	if exist {
		fmt.Println("删除存在文件：", "./"+pak_name+".pak")
		os.Remove("./" + pak_name + ".pak")
	}
	pakFile, err := os.OpenFile("./"+pak_name+".pak", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	defer pakFile.Close()
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)
	root := doc.CreateElement("map")
	var i = 0
	for rows.Next() {
		i++
		var Type int64
		var X int64
		var Y int64
		var Zoom int64
		var Tile []byte
		err := rows.Scan(&Type, &Zoom, &X, &Y, &Tile)
		if err != nil {
			fmt.Printf("get user info error [%s]", err)
		}
		n, _ := pakFile.Seek(0, 2)
		var znode, xnode, ynode *etree.Element
		znodelist := root.SelectElement(strings.Join([]string{"z", strconv.FormatInt(Zoom, 10)}, ""))
		if znodelist != nil {
			znode = znodelist
		} else {
			znode = root.CreateElement(strings.Join([]string{"z", strconv.FormatInt(Zoom, 10)}, ""))
		}
		xnodelist := znode.SelectElement(strings.Join([]string{"x", strconv.FormatInt(X, 10)}, ""))
		if xnodelist != nil {
			xnode = xnodelist
		} else {
			xnode = znode.CreateElement(strings.Join([]string{"x", strconv.FormatInt(X, 10)}, ""))
		}
		ynodelist := xnode.SelectElement(strings.Join([]string{"y", strconv.FormatInt(Y, 10)}, ""))
		if ynodelist != nil {
			ynode = ynodelist
		} else {
			ynode = xnode.CreateElement(strings.Join([]string{"y", strconv.FormatInt(Y, 10)}, ""))
		}
		ynode.CreateAttr("offset", strconv.FormatInt(n, 10))
		ynode.CreateAttr("length", strconv.Itoa(len(Tile)))
		pakFile.Write(Tile)
		if i%1000 == 0 {
			fmt.Println("当前时间 ：", string(time.Now().Format("2006-01-02 15:04:05")), "; 打包进度：", i, "; 级别:", Zoom)
		}
	}
	fmt.Println("打包完成总瓦块数：", i)
	existidx, err := PathExists("./" + pak_name + ".idx")
	if existidx {
		fmt.Println("删除存在文件：", "./"+pak_name+".idx")
		os.Remove("./" + pak_name + ".idx")
	}
	idxFile, err := os.OpenFile("./"+pak_name+".idx", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	defer idxFile.Close()
	doc.Indent(2)
	doc.WriteTo(idxFile)
	elapsed := time.Since(start)
	fmt.Println("地图打包执行完成!耗时：", elapsed)
	time.Sleep(time.Second * 2)
}

/*
如果返回的错误为nil,说明文件或文件夹存在
如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
如果返回的错误为其它类型,则不确定是否在存在
*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// 存在
		return true, nil
	}
	if os.IsNotExist(err) {
		// 不存在
		return false, nil
	}
	// 不存在
	return false, err
}
