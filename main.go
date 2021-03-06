package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	validator "gopkg.in/go-playground/validator.v8"
)

const authMethod = "JWT"
const apiVersion = "1.0.0"
const dbUsername = "root"
const dbPassword = ""
const dbDatabase = ""

var test int = func() int {
	return 1
}()

const secretKey string = "mocked-key"

var blacklisted_ip map[string]int
var lock sync.Mutex

func main() {

	//Create container

	Container.Bind("database", func(c *container) interface{} {
		db, err := gorm.Open("mysql", dbUsername+":"+dbPassword+"@/"+dbDatabase+"?charset=utf8&parseTime=True&loc=Local")

		if err != nil {
			panic(err)
		}

		return db
	})

	Container.Bind("user_provider", func(c *container) interface{} {
		return &MysqlUserProviderStruct{c.Resolve("database").(*gorm.DB)}
	})

	Container.Bind("guard", func(c *container) interface{} {
		return &JwtGuard{
			secret:       []byte(secretKey),
			userProvider: c.Resolve("user_provider").(UserProviderInterface),
			resolveBy:    "username",
		}
	})

	//Stop app if we don't have DB connection

	db := Container.Resolve("database").(*gorm.DB)

	db.AutoMigrate(&User{}, &Profile{})
	db.Model(&Profile{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	// var jwtGuard GuardInterface = &JwtGuard{
	// 	secret:       []byte(secretKey),
	// 	userProvider: mysqlUserProvider,
	// 	resolveBy:    "username",
	// }

	//jwtGuard.SetUserProvider(mysqlUserProvider)

	//Init variables

	//Create random seed for every random actions
	rand.Seed(time.Now().UTC().UnixNano())
	blacklisted_ip = make(map[string]int)
	lock = sync.Mutex{}

	//Bind custom validators

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("allowed_type", allowedType)
	}

	// Load router
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.Use(throttleMiddleware(10))
		v1.GET("/", indexAction)
		v1.POST("/login", loginAction)
		v1.POST("/register", registerAction)
	}

	v1auth := v1.Group("/")
	{
		v1auth.Use(authMiddleware(Container.Resolve("guard").(GuardInterface)))
		v1auth.GET("/user", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user": c.MustGet("user"),
			})
		})
	}
	router.Run()
}
