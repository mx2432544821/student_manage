package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

//student结构体
type student struct {
	ID      int
	Name    string
	Age     int
	Address string
	Class   int
	Avatar  string
}

//全局变量
var arr []student
var i int

//初始化数据
func Init() {
	arr = make([]student, 0, 50)
	for i = 1; i <= 40; i++ {
		switch rand.Int() % 7 {
		case 0:
			arr = append(arr, student{i, "minxu", 21, "beijing", 111183, "/img/1.jpg"})
		case 1:
			arr = append(arr, student{i, "hyf", 15, "wuhan", 111181, "/img/2.jpg"})
		case 2:
			arr = append(arr, student{i, "jjz", 14, "sichuan", 111183, "/img/3.jpg"})
		case 3:
			arr = append(arr, student{i, "fsz", 18, "henan", 111182, "/img/4.jpg"})
		case 4:
			arr = append(arr, student{i, "fyp", 26, "jaingxi", 111183, "/img/5.jpg"})
		case 5:
			arr = append(arr, student{i, "hy", 21, "yichang", 111181, "/img/6.jpg"})
		case 6:
			arr = append(arr, student{i, "ydy", 25, "shanghai", 111182, "/img/7.jpg"})
		}
	}
}

func main() {
	Init() //初始化数据
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("./templates/*")

	r.GET("/", index)
	r.POST("/update", update)
	r.GET("/delete", delStudent)
	r.POST("/add", add)
	r.Run(":8080")
}

func update(c *gin.Context) {
	pageNum := c.PostForm("pageNum")
	fmt.Println(pageNum,"update")
	var stu student
	//如果传入stu  会被拷贝 所以需要传递 指针/地址
	err := c.ShouldBind(&stu) //无论是get post（json）  都可以绑定数据
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	k:=0
	for ;k<len(arr);k++{
		if(arr[k].ID==stu.ID){
			stu.Avatar=arr[k].Avatar
			break
		}
	}
	file,err:=c.FormFile("file")
	if(err!=nil){
		fmt.Println("aa")
	}else{
		err:=c.SaveUploadedFile(file,"./static/img/"+file.Filename)
		if(err!=nil){
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		stu.Avatar="/img/"+file.Filename
		fmt.Println("new Avatar",stu.Avatar)
	}

	arr[k]=stu;
	c.Redirect(302,"/?pageNum="+pageNum)
	//c.JSON(200,gin.H{
	//	"message":stu,
	//})
}

func add(c *gin.Context) {
	var stu student
	//如果传入stu  会被拷贝 所以需要传递 指针/地址
	err := c.ShouldBind(&stu) //无论是get post（json）  都可以绑定数据
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	stu.ID=i
	i++

	file,err:=c.FormFile("file")
	fmt.Println(file,err)
	if(err!=nil){
		fmt.Println("FormFile failed")
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}else{
		err:=c.SaveUploadedFile(file,"./static/img/"+file.Filename)
		if(err!=nil){
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		stu.Avatar="/img/"+file.Filename
	}

	pages := len(arr) / 7
	if len(arr)%7 == 0 {
		pages--
	}
	arr = append(arr, stu)
	c.Redirect(302,"/?pageNum="+strconv.FormatInt(int64(pages), 10))
	//c.JSON(200, gin.H{
	//	"user": stu,
	//})
}

func delStudent(c *gin.Context) {
	strid := c.Query("id")
	pageNum := c.Query("pageNum")
	fmt.Println("delete", strid)
	id, err := strconv.ParseInt(strid, 10, 0)
	if err != nil {
		c.JSON(500, gin.H{"message": "删除参数不正确"})
		return
	}
	var index int
	for i := 0; i < len(arr); i++ {
		if arr[i].ID == int(id) {
			index = i
			break
		}
	}
	//删除索引为index的元素
	arr = append(arr[:index], arr[index+1:]...)
	c.Redirect(302, "/?pageNum="+pageNum)
	//c.JSON(200,gin.H{"message":index})
}

func index(c *gin.Context) {
	pageNumStr := c.DefaultQuery("pageNum", "0")
	pageNum, err := strconv.ParseInt(pageNumStr, 10, 0)
	if err != nil {
		c.JSON(500, gin.H{"message": "分页参数不正确"})
	}
	tag := c.Query("tag")
	fmt.Println("tag is", tag, pageNum)
	if tag == "pre" {
		c.Redirect(302, "/?pageNum="+strconv.FormatInt(pageNum-1, 10))
		return
	} else if tag == "post" {
		c.Redirect(302, "/?pageNum="+strconv.FormatInt(pageNum+1, 10))
		return
	}
	fmt.Println("tag is kong")
	if pageNum < 0 {
		c.Redirect(302, "/?pageNum=0")
		return
	}
	end := pageNum*7 + 7
	if end > int64(len(arr)) {
		end = int64(len(arr))
	}
	if int(pageNum*7) > len(arr)-1 {
		if int((pageNum-1)*7) > len(arr)-1 {
			c.JSON(400, gin.H{"message": "页数过大"})
		} else {
			pageNum--
			c.Redirect(302, "/?pageNum="+strconv.Itoa(int(pageNum)))
		}

		return
	}
	page := arr[pageNum*7 : end]
	pages := len(arr) / 7
	if len(arr)%7 == 0 {
		pages--
	}

	fmt.Println(page)
	c.HTML(200, "index.html", gin.H{
		"students": page,
		"pageNum":  pageNum,
		"lastNum":  pages,
	})
}
