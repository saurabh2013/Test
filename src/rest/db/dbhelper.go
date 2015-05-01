package db

import (
	"fmt"
	"github.com/gocql/gocql"
	"rest/log"
)

type dbClient struct {
	Session *gocql.Session
}
type Contact struct {
	Id    int
	Name  string
	Phone string
}

const (
	DB_SERVER = "localhost"
	KEYSPACE  = "dev"
)

//create table query
//CREATE TABLE contacts (id int PRIMARY KEY,name varchar, phone varchar)

func NewdbClient() *dbClient {
	cluster := gocql.NewCluster(DB_SERVER)
	cluster.Keyspace = KEYSPACE
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()
	return &dbClient{Session: session}
}

func (this *dbClient) DeleteUser(ct *Contact) error {
	if this.Session == nil {
		return fmt.Errorf("DB Connection Lost.")
	}
	if ct.Id < 1 {
		return fmt.Errorf("Please provide 'Contact.Id' input param")
	}
	err := this.Session.Query(fmt.Sprintf("DELETE FROM CONTACTS WHERE id = %d", ct.Id)).Exec()
	if err != nil {
		log.Info("ERROR: ", err)
	}
	defer this.Session.Close()
	return err
}

func (this *dbClient) InsertUser(ct *Contact) error {
	if this.Session == nil {
		return fmt.Errorf("DB Connection Lost.")
	}
	err := this.Session.Query(fmt.Sprintf(`INSERT INTO CONTACTS (id, name, phone) VALUES (%d,'%s','%s')`, ct.Id, ct.Name, ct.Phone)).Exec()
	if err != nil {
		log.Info("ERROR: ", err)
	}
	defer this.Session.Close()
	return err
}

func (this *dbClient) GetUsers(ct *Contact) ([]*Contact, error) {
	defer this.Session.Close()
	if this.Session == nil {
		return nil, fmt.Errorf("DB Connection Lost.")
	}
	q := "SELECT * FROM CONTACTS"
	if ct.Id > 0 {
		q = fmt.Sprintf("SELECT * FROM CONTACTS WHERE id=%d", ct.Id)
	}

	var name, phone string
	var id int

	var contacts []*Contact
	iter := this.Session.Query(q).Iter()
	for iter.Scan(&id, &name, &phone) {
		c := &Contact{Id: id, Name: name, Phone: phone}
		contacts = append(contacts, c)
	}
	err := iter.Close()
	if err != nil {
		log.Info(" ERROE in GET e: ", err)
	}

	return contacts, err
}
