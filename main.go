package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"is-a-dev.com/autoapi/lib"

	"github.com/howeyc/gopass"
	_ "github.com/ziutek/mymysql/godrv"
)

const (
	// Ugly way to check to see if they passed in a password
	// chance of collision with a GUID is very low
	defaultPassValue = "5e7dc6f6a1a94c39be95b88a47c2458b"
)

var (
	dbPort  string
	dbHost  string
	dbName  string
	dbUname string
	dbPass  string
	verbs   string
)

func init() {
	flag.StringVar(&dbPort, "P", "3306", "database port")
	flag.StringVar(&dbPass, "p", defaultPassValue, "database password")
	flag.StringVar(&dbHost, "h", "localhost", "database host")
	flag.StringVar(&dbName, "d", "", "database name")
	flag.StringVar(&dbUname, "u", "root", "database username")
	flag.StringVar(&verbs, "v", "get", "HTTP verbs to handle.")
	flag.Parse()
}

func main() {
	pass := dbPass
	if pass == defaultPassValue {
		fmt.Print("Password:")
		pass = strings.TrimSpace(string(gopass.GetPasswdMasked()))
	}

	if strings.TrimSpace(dbPort) == "" {
		fmt.Println("Missing port")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbHost) == "" {
		fmt.Println("Missing host")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbName) == "" {
		fmt.Println("Missing database name")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbUname) == "" {
		fmt.Println("Missing username")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dbConn, err := sql.Open("mymysql", fmt.Sprintf("tcp:%s:%s*%s/%s/%s", dbHost, dbPort, dbName, dbUname, pass))
	if err != nil {
		flag.PrintDefaults()
		log.Panic(err)
	}
	verbs = strings.ToLower(verbs)
	err = lib.Generate(dbConn, dbName, verbs)
	if err != nil {
		panic(err)
	}
}
