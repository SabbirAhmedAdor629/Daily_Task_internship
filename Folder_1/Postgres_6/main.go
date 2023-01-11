package main

import (
	"database/sql"
	"fmt"

	//	"time"
	//	"log"
	//	"time"
	_ "github.com/lib/pq"
)

var db *sql.DB



const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "messaging_development"
)



// MESSAGE TEMPLATE TABLE
type MessageTemplate struct {
	Active     bool
	DailyCap   int
	MonthlyCap int
	ToatalCap  int
}
func ReadFromMessageTempaltes(db *sql.DB) {
	sqlStatement := `SELECT active, daily_cap, monthly_cap, total_cap FROM message_templates WHERE id=$1`

	var Template MessageTemplate
	err := db.QueryRow(sqlStatement, 1).Scan(&Template.Active, &Template.DailyCap, &Template.MonthlyCap, &Template.ToatalCap)
	if err != nil {
		panic(err)
	}
}
func (mt *MessageTemplate) messageTemplateActive() bool {
	return mt.Active
}



// PUSH-REGISTRATION TABLE
type PushRegistration struct {
	AwsArnStatus string
	Sandbox      bool
}
func ReadFromPushRegistration(db *sql.DB) {
	sqlStatement := `SELECT aws_arn_status, sandbox FROM push_registrations WHERE id=$1`

	var PushRegistration PushRegistration
	err := db.QueryRow(sqlStatement, 1).Scan(&PushRegistration.AwsArnStatus, &PushRegistration.Sandbox)
	if err != nil {
		panic(err)
	}
	fmt.Println(PushRegistration.AwsArnStatus)
	fmt.Println(PushRegistration.Sandbox)
}
func (pr *PushRegistration) active() bool {
	return pr.AwsArnStatus == "registered" && !pr.Sandbox
}



// CAMPAIGN TABLE
type CampaignTable struct {
	ScheduleTime      int
	EventName         string
	MessageTemplateId int
}
func ReadFromCampaignTable(db *sql.DB) {
	sqlStatement := `SELECT schedule_time, event_name, message_template_id FROM campaigns WHERE id=$1`
	var CampaignTable CampaignTable
	err := db.QueryRow(sqlStatement, 4).Scan(&CampaignTable.ScheduleTime, &CampaignTable.EventName, &CampaignTable.MessageTemplateId)
	if err != nil {
		panic(err)
	}
}
func (ct *CampaignTable) schedule_time() int {
	return ct.ScheduleTime
}


// MAIN FUNCTION
func main() {
	// DB connection
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	// ReadFromPushRegistration(db)
	// pr := &PushRegistration{}
	// fmt.Println(pr.active())

	// ReadFromMessageTempaltes(db)
	// mt := &MessageTemplate{Active: true}
	// fmt.Println(mt.messageTemplateActive())

	ReadFromCampaignTable(db)
	ct := &CampaignTable{}
	fmt.Println(ct.schedule_time())

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
