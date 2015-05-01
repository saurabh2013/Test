package handler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"rest/consts"
	//"runtime/debug"
	"rest/db"
	"rest/log"
	"strings"
)

type Contacts struct {
	XMLName xml.Name  `xml:"contacts"`
	C       []Contact `xml:"contact"`
}
type Contact struct {
	Id    int    `xml:"id"`
	Name  string `xml:"name,omitempty"`
	Phone string `xml:"phone,omitempty"`
}

func ProcessError(w http.ResponseWriter, r *http.Request, e error) {
	if e != nil {
		fmt.Fprintf(w, "Error: %s", e.Error())
		//fmt.Printf("STACK: %s", string(debug.Stack()))
	}
}

func ProcessResponse(w http.ResponseWriter, r *http.Request, d []*db.Contact) {
	if d == nil {
		ProcessError(w, r, fmt.Errorf("Can not find data"))
	}
	if aHdr := r.Header.Get(consts.HEADER_ACCEPT); strings.EqualFold(aHdr, consts.APPLICATION_XML) {
		parseXML(w, r, d)
	} else {
		parseJSON(w, r, d)
	}
}
func parseXML(w http.ResponseWriter, r *http.Request, d []*db.Contact) {

	w.Header().Set(consts.HEADER_CONTENTTYPE, consts.APPLICATION_XML)

	var cnts Contacts
	for _, v := range d {
		cnts.C = append(cnts.C, Contact{Id: v.Id, Name: v.Name, Phone: v.Phone})
	}
	//log.Info("RESPONSE XML PRS ", cnts.C)
	buffer, er := xml.MarshalIndent(&cnts, "", "    ")
	if er != nil {
		ProcessError(w, r, er)
		return
	}
	xmlContent := xml.Header + bytes.NewBuffer(buffer).String()
	fmt.Fprint(w, xmlContent)

}
func parseJSON(w http.ResponseWriter, r *http.Request, d []*db.Contact) {
	if len(d) < 1 {
		ProcessError(w, r, fmt.Errorf("Empty response item."))
		return
	}
	w.Header().Set(consts.HEADER_CONTENTTYPE, consts.APPLICATION_JSON)
	res, err := json.MarshalIndent(d, "", "  ")

	if err != nil {
		ProcessError(w, r, err)
		return
	}
	fmt.Fprint(w, string(res))
}
