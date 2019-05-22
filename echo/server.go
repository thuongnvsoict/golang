package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
)

var db *sql.DB

type User struct{
	Id int `json:"id" form:"id" query:"id"`
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
	PhoneNumber	string `json:"phonenumber" form:"phonenumber" query:"phonenumber"`
	Email string `json:"email" form:"email" query:"email"`
}

type Cat struct{
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct{
	Name string `json:"name"`
	Type string `json:"type"`
}

//Path Param

//Query Param
func show(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")

	if dataType == "string"{
		return c.String(http.StatusOK, "Name:" + catName + " , Type:" + catType)
	}

	if dataType == "json"{
		return c.JSON(http.StatusOK, map[string]string{
			"Name" : catName,
			"Type" : catType,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error" : "You need to let us know data type that you want to get",
	})

}

func addCat(c echo.Context) error{
	cat := Cat{}
	//fmt.Println("xxxxxx")
	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)

	if err != nil{
		log.Printf("Failed reading the request body %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)

	if err != nil{
		log.Printf("Failed unmarshaling in addCats %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("This is your cat %#v", cat)
	return c.String(http.StatusOK, "We got your cat")

}

func getUser(c echo.Context) error{

	fmt.Println(db)
	username := c.QueryParam("username")
	result, err := db.Query("SELECT * FROM users WHERE username = ?" , username)

	if err != nil {
		panic(err.Error())
	}

	var user User
	if result.Next(){
		err = result.Scan(&user.Id, &user.Username, &user.Password, &user.PhoneNumber, &user.Email)
		if err != nil {
			panic(err.Error())
		}

	}
	return c.JSON(http.StatusOK, user)
	// return c.JSON(http.StatusOK, "Lay duoc roi nay")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Horay you are on secret admin amin page!")
}

func insertUser(c echo.Context) error{
	insert, err := db.Query("insert into users(`username`, `password`, `phonenumber`, `email`) value ('tunglx', '12345678', '0365681245', 'tunglx@gmail.com')")
	if err != nil{
		panic(err.Error())
	}
	defer insert.Close()
	return c.String(http.StatusOK, "We got new user")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH,    echo.POST, echo.DELETE},
	}))
	/*------------
		MYSQL
	 ------------*/
	fmt.Println("Go MySql Tutorial")

	//Connect
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/account")
	if err != nil {
		panic(err.Error())
	}else{
		fmt.Println("Connected to Mysql Server!")
	}
	fmt.Println(db)
	defer db.Close()

	// Insert



	/*------------
		ROUTING
	 ------------*/
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!!!")
	})

	g:= e.Group("/admin")
	g.GET("/main", mainAdmin)

	e.GET("/users", getUser)

	e.POST("/users", insertUser)

	/*------------
		START
	 ------------*/
	e.Logger.Fatal(e.Start(":1323"))


}
