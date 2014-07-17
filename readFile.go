package main

import(
	"os"
	"encoding/csv"
	_ "github.com/lib/pq"
	"database/sql"
	"path/filepath"
	"strings"
)

func check(e error){
	if e!=nil{
		panic(e)
	}
}

func main(){
	file, err := os.Open("./routes.txt")
	check(err)
	db, err := sql.Open("postgres", "user=postgres dbname=postgres host=localhost password=ohlongjohnson port=5432 sslmode=disable")
	defer db.Close()
	check(err)
	var filename = file.Name()
	var extension = filepath.Ext(filename)
	var cleanName = filename[2:len(filename)-len(extension)]
	var str = []string{"CREATE TABLE ",cleanName,"()"}
	_, err = db.Exec(strings.Join(str,""))
	check(err)
	reader := csv.NewReader(file)
	record, err := reader.Read()
	valuesa := strings.Join(record,"\",\"")
	values := strings.Join([]string{"(\"",valuesa,"\")"},"")
	check(err)
	for _, col := range record{
		var str = []string{"ALTER TABLE ",cleanName," ADD ",col," text"}
		db.Exec(strings.Join(str,""))
	}
	check(err)
	for record, err = reader.Read();record!=nil;record,err=reader.Read(){
		check(err)
		for i:=0;i<len(record);i++{
			if strings.EqualFold(record[i],""){
				record[i]=" "
			}
		}
		args := []string{"INSERT INTO ",cleanName," ",values," VALUES ('",strings.Join(record,"','"),"')"}
		db.Exec(strings.Join(args,""))
	}
}
