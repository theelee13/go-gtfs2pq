package main

import(
	"os"
	"encoding/csv"
	_ "github.com/lib/pq"
	"database/sql"
	"path/filepath"
	"strings"
	"fmt"
)

func check(e error){
	if e!=nil{
		panic(e)
	}
}

func main(){
	db, err := sql.Open("postgres", "user=postgres dbname=postgres host=localhost password=ohlongjohnson port=5432 sslmode=disable")
	defer db.Close()
	check(err)
	files := []string{"./agency.txt","./routes.txt","./calendar.txt","./calendar_dates.txt","./stops.txt","./trips.txt"}
	for _, filename  := range files{
		file, err := os.Open(filename)
		if err!=nil{
			fmt.Printf("error loading data for %s\n",filename)
			continue
		}
		extension := filepath.Ext(filename)
		cleanName := filename[2:len(filename)-len(extension)]
		str := []string{"CREATE TABLE ",cleanName,"()"}
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
		fmt.Printf("finished loading data for %s\n",filename)
	}
}
