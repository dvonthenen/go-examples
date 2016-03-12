package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/icza/gowut/gwu"
)

type account struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type accounts []account

type myButtonHandler struct {
	usernames []gwu.TextBox
	names     []gwu.TextBox
	emails    []gwu.TextBox
}

func (h *myButtonHandler) HandleEvent(e gwu.Event) {
	if b, isButton := e.Src().(gwu.Button); isButton {
		for i := 0; i < 3; i++ {
			fmt.Println("Username: " + h.usernames[i].Text())
			fmt.Println("Name: " + h.names[i].Text())
			fmt.Println("Email: " + h.emails[i].Text())

			if len(h.usernames[i].Text()) == 0 ||
				len(h.names[i].Text()) == 0 ||
				len(h.emails[i].Text()) == 0 {
				continue
			}

			//TODO Call Add REST API
		}
		e.MarkDirty(b)
	}
}

func main() {
	// Create and build a window
	win := gwu.NewWindow("main", "3-Tier App Demo")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	//Add users
	addform := gwu.NewVerticalPanel()
	addform.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	addform.SetCellPadding(10)
	addform.SetAttr("width", "400")
	addform.Add(gwu.NewLabel("Add New Accounts"))

	var usernames []gwu.TextBox
	var names []gwu.TextBox
	var emails []gwu.TextBox

	for i := 0; i < 3; i++ {
		p := gwu.NewHorizontalPanel()
		p.SetCellPadding(2)
		p.Add(gwu.NewLabel("Username:"))
		tbusername := gwu.NewTextBox("")
		p.Add(tbusername)
		p.Add(gwu.NewLabel("Name:"))
		tbname := gwu.NewTextBox("")
		p.Add(tbname)
		p.Add(gwu.NewLabel("Email:"))
		tbemail := gwu.NewTextBox("")
		p.Add(tbemail)

		usernames = append(usernames, tbusername)
		names = append(names, tbname)
		emails = append(emails, tbemail)

		addform.Add(p)
	}

	btnadd := gwu.NewButton("Add")
	btnadd.SetAttr("align", "center")
	btnadd.AddEHandler(&myButtonHandler{usernames, names, emails}, gwu.ETypeClick)
	addform.Add(btnadd)

	win.Add(addform)

	//Display users...
	acctlist := gwu.NewVerticalPanel()
	acctlist.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	acctlist.SetCellPadding(10)
	acctlist.SetAttr("width", "400")

	p := gwu.NewHorizontalPanel()
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("Current List of Accounts"))
	btnrefresh := gwu.NewButton("Refresh")
	btnrefresh.SetAttr("align", "center")
	//btnrefresh.AddEHandler(&myButtonHandler{usernames, names, emails}, gwu.ETypeClick)
	p.Add(btnrefresh)

	acctlist.Add(p)

	//TODO: REST API call to GET... loop adding labels
	url := "http://127.0.0.1:9000/user"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	bodystr, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodystr))

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	var accts []account
	err = json.Unmarshal(body, &accts)
	/*
		for acct := range accts {
			fmt.Printf("Id: %v", acct.Id)
			fmt.Println("Username: %s", acct.Username)
			fmt.Println("Name: %s", acct.Name)
			fmt.Println("Email: %s", acct.Email)
		}
	*/
	//*
	for i := 0; i < len(accts); i++ {
		fmt.Printf("Id: %v", accts[i].Id)
		fmt.Println("Username: %s", accts[i].Username)
		fmt.Println("Name: %s", accts[i].Name)
		fmt.Println("Email: %s", accts[i].Email)
	}
	//	*/

	win.Add(acctlist)

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("", "localhost:8081")
	server.SetText("Test GUI App")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser
}
