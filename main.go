package main

// Import part
import (
	"database/sql"
	"fmt"
	"github.com/gorilla/schema"
	_ "log"
	"net/http"
	_ "net/http"
	"os"
	_ "strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq" // here
)

//Declare part

// -- Variable --
var db *sql.DB

// ** Dùng để decode request param
var decoder = schema.NewDecoder()

// -- Constant --
const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

//DB Connection part
func initDB() {
	config := dbConfig()
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	conf[dbhost] = "ec2-107-22-241-205.compute-1.amazonaws.com"
	conf[dbport] = "5432"
	conf[dbuser] = "vvmtprwhttnfzu"
	conf[dbpass] = "f3431810665873fb90038634e1b0272639ac8436fa63d8fff03039ffb5f4e213"
	conf[dbname] = "dbvo3qeb8d6d9l"

	return conf
}

// Main Function
func main() {
	// Init DB connection
	//initDB()
	//defer db.Close()

	// Init route for rq
	router := gin.New()
	router.LoadHTMLFiles("index.html")

	// Main Screen
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			},
		)
	})

	// CMD part
	/*router.GET("/cmd", func(c *gin.Context) {
		data := c.Request.URL.Query()
		rs := data.Get("data")
		var cmdRepo cmdRepo
		json.Unmarshal([]byte(rs), &cmdRepo)
		fmt.Println(rs)
		c.JSON(200, gin.H{"msg": cmdRepo})
	})*/

	// Chat Screen
	router.GET("/main/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Global Chat
	router.GET("/main", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Socket Request
	router.GET("/ws/:roomId/*userId", func(c *gin.Context) {
		roomID := c.Param("roomId")
		userIDRaw := c.Param("userId")
		rp := strings.Replace(userIDRaw, "/", "", 1)
		if roomID == "" {
			roomID = "1"
		}
		// Init Socket Service
		serveWs(c.Writer, c.Request, roomID, rp)
	})

	// Declare port running on server
	var port = ":" + os.Getenv("PORT")
	_ = router.Run(port)

	// For local test purpose
	//_ = router.Run(":8080")
}

// All Function Here

// DB Function
