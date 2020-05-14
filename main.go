package main

import (
	"MapPackageGo/mapinit"
	"fmt"
	"github.com/beevik/etree"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now() // 获取当前时间

	mapinit.MapInit()
	pak_name:=mapinit.Config["pak_name"]
	map_type:=mapinit.Config["map_type"]
	zoom_min:=mapinit.Config["zoom_min"]
	zoom_max:=mapinit.Config["zoom_max"]
	gmapnetcache:=mapinit.Config["tablename"]
	 sql := " select Type, Zoom, X, Y, Tile from "+gmapnetcache+"  where Type = "+map_type+" and Zoom >= "+zoom_min+" and Zoom <= "+zoom_max+" ORDER BY zoom,x,y"

	rows,err := mapinit.Db.Query(sql)
	if err != nil{
		fmt.Printf("select fail [%s]",err)
	}
	pakFile, err := os.OpenFile("./"+pak_name+".pak", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	doc:= etree.NewDocument()
	root:= doc.CreateElement("map")
	var i=0
	for rows.Next(){
		i++
		var Type int64
		var X int64
		var Y int64
		var Zoom int64
		var Tile []byte
		err := rows.Scan(&Type,&Zoom,&X,&Y,&Tile)
		if err != nil{
			fmt.Printf("get user info error [%s]",err)
		}
		n, _ := pakFile.Seek(0, 2)
		var znode,xnode,ynode *etree.Element
		znodelist:=root.SelectElement(strings.Join([]string{"z", strconv.FormatInt(Zoom,10)},""))
		if znodelist!=nil{
			znode =znodelist
		}else {
			znode=root.CreateElement(strings.Join([]string{"z", strconv.FormatInt(Zoom,10)},""))
		}
		xnodelist:=znode.SelectElement(strings.Join([]string{"x", strconv.FormatInt(X,10)},""))
		if xnodelist!=nil{
			xnode =xnodelist
		}else {
			xnode=znode.CreateElement(strings.Join([]string{"x", strconv.FormatInt(X,10)},""))
		}
		ynodelist:=xnode.SelectElement(strings.Join([]string{"y", strconv.FormatInt(Y,10)},""))
		if ynodelist!=nil{
			ynode =ynodelist
		}else {
			ynode=xnode.CreateElement(strings.Join([]string{"y", strconv.FormatInt(Y,10)},""))
		}
		ynode.CreateAttr("offset",strconv.FormatInt(n,10))
		ynode.CreateAttr("length",strconv.Itoa(len(Tile)))
		pakFile.Write(Tile)
		/*fmt.Println("追加文件长度：%i",len(Tile))
		fmt.Println("文件当前总长度：%i",n)
		fmt.Println(Zoom)
		fmt.Println(X)
		fmt.Println(Y)*/
		if i%10000==0{
			fmt.Println("打包进度：",i )
		}
	}
	fmt.Println("打包完成：",i )
	idxFile, err := os.OpenFile("./"+pak_name+".idx", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	doc.Indent(2)
	doc.WriteTo(idxFile)
	elapsed := time.Since(start)
	fmt.Println("地图打包执行完！成耗时：秒", elapsed)

	time.Sleep(time.Second*2)
}