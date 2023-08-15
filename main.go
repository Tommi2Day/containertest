package main

import (
	"database/sql"
	"fmt"
	"github.com/tommi2day/gomodules/dblib"
	go_ora "github.com/sijms/go-ora/v2"
	"os"
	"path"
	"strings"
)

var tnsEntries dblib.TNSEntries
var domain string

func main() {
	urlOptions := map[string]string{
		"trace file": "trace.log",
		"timeout":    "5",
	}
	args := os.Args[1:]
	if len(args) < 2 {
		println("Usage: $0 <CDB> <PDB>")
		return
	}
	dbuser := os.Getenv("DBUSER")	
	dbpass := os.Getenv("DBPASS")
	cdb := args[0]
	pdb := args[1]
	err := loadTNS()
	if err != nil {
		fmt.Println(err)
		return
	}
	tnsDesc, err := getTNSDesc(cdb)
	if tnsDesc == "" {
		fmt.Println("Service not found")
		return
	}
	databaseURL := go_ora.BuildJDBC(dbuser, dbpass, tnsDesc, urlOptions)
	conn, err := sql.Open("oracle", databaseURL)
	if err != nil {
		fmt.Printf("can't connect to %s:%s\n", cdb, err)
		return
	}
	_, err = conn.Exec("alter session set container=" + pdb)
	if err != nil {
		fmt.Printf("can't set container to %s:%s\n", pdb, err)
		return
	}
	sql := "select sys_context('USERENV','CON_NAME') from dual"
	row := conn.QueryRow(sql)
	if row == nil {
		fmt.Println("Row not returned")
		return
	}
	var conName string
	err = row.Scan(&conName)
	if err != nil {
		fmt.Printf("Query solut not expected:%s", err)
		return
	}
	if strings.ToUpper(pdb) != conName {
		fmt.Printf("PDB to not match, exp:%s, actual %s\n", pdb, conName)
		return
	}
	fmt.Println("Test OK")
}

func loadTNS() (err error) {
	tnsadmin := os.Getenv("TNS_ADMIN")
	tnsFile := path.Join(tnsadmin, "tnsnames.ora")
	// load available tns entries
	tnsEntries, domain, err = dblib.GetTnsnames(tnsFile, true)
	l := len(tnsEntries)
	if err != nil {
		return
	}
	if l == 0 {
		err = fmt.Errorf("cannot proceed without tns entries")

	}
	return
}

func getTNSDesc(dbservice string) (tnsDesc string, err error) {
	if dbservice == "" {
		err = fmt.Errorf("dont have a service to lookup")
		return
	}

	fmt.Printf("get info for service %s ", dbservice)
	entry, found := dblib.GetEntry(dbservice, tnsEntries, domain)
	if !found {
		err = fmt.Errorf("alias %s not found", dbservice)
		return
	}

	desc := entry.Desc
	repl := strings.NewReplacer("\r", "", "\n", "", "\t", "", " ", "")
	tnsDesc = repl.Replace(desc)
	return
}
