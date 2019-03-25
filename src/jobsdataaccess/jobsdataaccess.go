package jobsdataaccess

import "database/sql"

type JobsDataAccess struct {
	Db sql.DB
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
