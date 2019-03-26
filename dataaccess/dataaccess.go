
package dataaccess

import "database/sql"
import "log"
type JobsDataAccess struct {
	Db sql.DB
}

func (d *JobsDataAccess) Init(){
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password=postgres dbname=karriereat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	d.Db = *db
}

func (d *JobsDataAccess) SaveJobDetails(jobDetails *JobsDetails) {
	_, err := d.Db.Exec(`INSERT INTO jobs1(url, title, company,location, date, content)
	VALUES($1,$2,$3,$4,$5,$6)`, jobDetails.Url, jobDetails.Title, jobDetails.Company, jobDetails.Location, jobDetails.Date, jobDetails.Content)
	if err != nil {
		panic(err)
	}
}

type JobsDetails struct {
	Url      string
	Title    string
	Company  string
	Location string
	Date     string
	Content  string
}


//https://linuxhint.com/install-pgadmin4-ubuntu/