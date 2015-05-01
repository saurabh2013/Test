package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"rest/core/util"
	"rest/db"
	//"rest/log"
)

type User struct {
	Id    int
	Name  string
	Phone string
}

const (
	STS = "Succssesfull DB Transaction for '%s'"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if usr, err := parseGetReq(r); err != nil {
			ProcessError(w, r, err)
		} else {
			usr.GetUserDetails(w, r)
		}
	} else if r.Method == "POST" || r.Method == "PUT" {
		if usr, err := parsePostReq(r); err != nil {
			ProcessError(w, r, err)
		} else {
			usr.PostUserDetails(w, r)
		}
	} else if r.Method == "DELETE" {
		if usr, err := parsePostReq(r); err != nil {
			ProcessError(w, r, err)
		} else {
			usr.DeleteUserDetails(w, r)
		}
	}
}

func parseGetReq(r *http.Request) (*User, error) {
	usr := &User{}
	u, err := url.Parse(r.RequestURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse url: %s", r.RequestURI)
	}
	//parse query
	q, err := url.ParseQuery(u.RawQuery)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse query in url: %s", r.RequestURI)
	}
	if value := q.Get("Id"); value != "" {

		if _id := util.GetIntValue(value); _id != nil {
			usr.Id = *_id
		} else {
			return nil, fmt.Errorf("Failed to Get Id param in query url: %s", r.RequestURI)
		}
	}

	return usr, nil
}
func parsePostReq(r *http.Request) (*User, error) {
	var dataJson string
	if b, er := ioutil.ReadAll(r.Body); er != nil {
		return nil, er
	} else {
		dataJson = string(b)
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dataJson), &data); err != nil {
		return nil, fmt.Errorf("Invalid input JSON Data, Failed to collect Expected Json")
	}
	if dbContact, err := GetContactInfo(data); err == nil {
		return dbContact, nil
	} else {
		return nil, err
	}
}

//Get Contact details
func (this *User) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	//DB operation for Get
	dbc := db.NewdbClient()
	resp, er := dbc.GetUsers(&db.Contact{Id: this.Id})

	if er != nil {
		ProcessError(w, r, er)
		return
	}
	//test data
	ProcessResponse(w, r, resp)
	//return nil, fmt.Errorf("Coud not get data for this id.")
}

func (this *User) PostUserDetails(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	//DB operation for Post
	dbc := db.NewdbClient()
	er := dbc.InsertUser(&db.Contact{Id: this.Id, Name: this.Name, Phone: this.Phone})
	if er != nil {
		ProcessError(w, r, er)
		return nil, nil
	}
	if er != nil {
		ProcessError(w, r, er)
		return nil, nil
	}

	fmt.Fprintf(w, fmt.Sprintf(STS, this.Name))
	return nil, nil
}

func (this *User) DeleteUserDetails(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	//DB operation for Post
	dbc := db.NewdbClient()
	er := dbc.DeleteUser(&db.Contact{Id: this.Id, Name: this.Name, Phone: this.Phone})
	if er != nil {
		ProcessError(w, r, er)
		return nil, nil
	}
	if er != nil {
		ProcessError(w, r, er)
		return nil, nil
	}

	fmt.Fprintf(w, fmt.Sprintf(STS, this.Name))
	return nil, nil
}

func GetContactInfo(data map[string]interface{}) (*User, error) {
	var bkc *User
	if len(data) > 0 {
		if _ct, k := data["Contact"]; k {
			if ct, k := _ct.(map[string]interface{}); k {
				bkc = &User{}
				if _id, k := ct["Id"]; k {
					if v := util.GetIntValue(_id); v == nil {
						return nil, fmt.Errorf("Invalid input 'Contact.Id' value type. Expecting integer..")
					} else {
						bkc.Id = *v
					}
				}
				if n, k := ct["Name"]; k {
					bkc.Name = n.(string)
				}
				if p, k := ct["Phone"]; k {
					bkc.Phone = p.(string)
				}
				if bkc != nil {
					return bkc, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Could Not able to parse input json request. Please verify Input request.")
}
