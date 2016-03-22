package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/icza/gowut/gwu"
)

type account struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type accounts []account

func refreshRestApi(address string, port int) accounts {
	url := "http://" + address + ":" + strconv.Itoa(port) + "/user"
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

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", string(body))

	var accts accounts
	err = json.Unmarshal(body, &accts)

	for i := 0; i < len(accts); i++ {
		fmt.Println("Id:", accts[i].Id)
		fmt.Println("Username:", accts[i].Username)
		fmt.Println("Name:", accts[i].Name)
		fmt.Println("Email:", accts[i].Email)
	}

	return accts
}

func addRestApi(address string, port int, accts accounts) accounts {
	url := "http://" + address + ":" + strconv.Itoa(port) + "/user"
	fmt.Println("URL:>", url)

	response, err := json.MarshalIndent(accts, "", "  ")
	if err != nil {
		panic(err) //not expecting error... just a short cut
	}

	fmt.Println("response:", string(response))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(response))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", string(body))

	var newaccts accounts
	err = json.Unmarshal(body, &newaccts)

	for i := 0; i < len(accts); i++ {
		fmt.Println("Id:", newaccts[i].Id)
		fmt.Println("Username:", newaccts[i].Username)
		fmt.Println("Name:", newaccts[i].Name)
		fmt.Println("Email:", newaccts[i].Email)
	}

	return newaccts
}

func deleteRestApi(address string, port int, id int) {
	url := "http://" + address + ":" + strconv.Itoa(port) + "/user/" + strconv.Itoa(id)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("DELETE", url, nil)
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", string(body))
}

type myBtnDelete struct {
	address string
	port    int
	id      int
	parent  gwu.Panel
	panel   gwu.Panel
}

func (h *myBtnDelete) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {

		fmt.Println("Delete called")
		fmt.Println("Id:", h.id)
		fmt.Println("Parent:", h.parent.Id().String())
		fmt.Println("Panel:", h.panel.Id().String())

		deleteRestApi(h.address, h.port, h.id)

		h.parent.Remove(h.panel)

		e.MarkDirty(h.parent)
	}
}

type myBtnAdd struct {
	address   string
	port      int
	usernames []gwu.TextBox
	names     []gwu.TextBox
	emails    []gwu.TextBox
	acctlist  gwu.Panel
}

func (h *myBtnAdd) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {

		fmt.Println("Add called")

		var accts accounts

		for i := 0; i < 3; i++ {
			fmt.Println("Username: " + h.usernames[i].Text())
			fmt.Println("Name: " + h.names[i].Text())
			fmt.Println("Email: " + h.emails[i].Text())

			if len(h.usernames[i].Text()) == 0 ||
				len(h.names[i].Text()) == 0 ||
				len(h.emails[i].Text()) == 0 {
				continue
			}

			accts = append(accts, account{0, h.usernames[i].Text(), h.names[i].Text(), h.emails[i].Text()})

			h.usernames[i].SetText("")
			h.names[i].SetText("")
			h.emails[i].SetText("")

			e.MarkDirty(h.usernames[i])
			e.MarkDirty(h.names[i])
			e.MarkDirty(h.emails[i])
		}

		newaccts := addRestApi(h.address, h.port, accts)

		for i := 0; i < len(newaccts); i++ {
			p := gwu.NewHorizontalPanel()
			p.SetCellPadding(2)
			p.Add(gwu.NewLabel("Username: " + newaccts[i].Username))
			p.Add(gwu.NewLabel("Name: " + newaccts[i].Name))
			p.Add(gwu.NewLabel("Email: " + newaccts[i].Email))

			btndelete := gwu.NewButton("Delete")
			btndelete.SetAttr("align", "center")
			btndelete.AddEHandler(&myBtnDelete{h.address, h.port, newaccts[i].Id, h.acctlist, p}, gwu.ETypeClick)
			p.Add(btndelete)

			fmt.Println("Panel Id:", p.Id().String())

			h.acctlist.Add(p)
		}

		e.MarkDirty(h.acctlist)
	}
}

func refresh(address string, port int, parent gwu.Panel) {
	accts := refreshRestApi(address, port)

	for i := 0; i < len(accts); i++ {
		p := gwu.NewHorizontalPanel()
		p.SetCellPadding(2)
		p.Add(gwu.NewLabel("Username: " + accts[i].Username))
		p.Add(gwu.NewLabel("Name: " + accts[i].Name))
		p.Add(gwu.NewLabel("Email: " + accts[i].Email))

		btndelete := gwu.NewButton("Delete")
		btndelete.SetAttr("align", "center")
		btndelete.AddEHandler(&myBtnDelete{address, port, accts[i].Id, parent, p}, gwu.ETypeClick)
		p.Add(btndelete)

		fmt.Println("Panel Id:", p.Id().String())

		parent.Add(p)
	}
}

func adduserform(address string, port int, acctlist gwu.Panel) gwu.Panel {
	addform := gwu.NewVerticalPanel()
	addform.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	addform.SetCellPadding(10)
	addform.SetAttr("width", "800")
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
	btnadd.AddEHandler(&myBtnAdd{address, port, usernames, names, emails, acctlist}, gwu.ETypeClick)
	addform.Add(btnadd)

	return addform
}

func listuserform(address string, port int) gwu.Panel {
	acctlist := gwu.NewVerticalPanel()
	acctlist.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	acctlist.SetCellPadding(10)
	acctlist.SetAttr("width", "800")

	acctlist.Add(gwu.NewLabel("Current List of Accounts"))

	//get current list of accounts
	refresh(address, port, acctlist)

	return acctlist
}

func main() {
	//define flags
	var port int
	flag.IntVar(&port, "port", 9000, "the REST server in which to bind to")
	var address string
	flag.StringVar(&address, "address", "127.0.0.1", "the REST server in which to bind to")
	//parse
	flag.Parse()

	// Create and build a window
	win := gwu.NewWindow("main", "3-Tier App Demo")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	//Display users...
	acctlist := listuserform(address, port)
	win.Add(acctlist)

	//Add users
	win.Add(adduserform(address, port, acctlist))

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("", "127.0.0.1:8000")
	server.SetText("Test GUI App")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser
}
