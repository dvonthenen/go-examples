package main

import (
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

func adduserform() gwu.Panel {
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
	btnadd.AddEHandler(&myBtnAdd{usernames, names, emails}, gwu.ETypeClick)
	addform.Add(btnadd)

	return addform
}

type myBtnDelete struct {
	id     int
	parent gwu.Panel
	panel  gwu.Panel
}

func (h *myBtnDelete) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {

		fmt.Println("Delete called")
		fmt.Println("Id:", h.id)
		fmt.Println("Parent:", h.parent.Id().String())
		fmt.Println("Panel:", h.panel.Id().String())

		//TODO Call Add REST API

		h.parent.Remove(h.panel)

		e.MarkDirty(h.parent)
	}
}

type myBtnAdd struct {
	usernames []gwu.TextBox
	names     []gwu.TextBox
	emails    []gwu.TextBox
}

func (h *myBtnAdd) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {

		fmt.Println("Add called")

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

			h.usernames[i].SetText("")
			h.names[i].SetText("")
			h.emails[i].SetText("")

			e.MarkDirty(h.usernames[i])
			e.MarkDirty(h.names[i])
			e.MarkDirty(h.emails[i])
		}
	}
}

func refresh(address string, port int, parent gwu.Panel, panels []gwu.Panel) {
	for i := 0; i < len(panels); i++ {
		parent.Remove(panels[i])
	}

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

		p := gwu.NewHorizontalPanel()
		p.SetCellPadding(2)
		p.Add(gwu.NewLabel("Username: " + accts[i].Username))
		p.Add(gwu.NewLabel("Name: " + accts[i].Name))
		p.Add(gwu.NewLabel("Email: " + accts[i].Email))

		btndelete := gwu.NewButton("Delete")
		btndelete.SetAttr("align", "center")
		btndelete.AddEHandler(&myBtnDelete{accts[i].Id, parent, p}, gwu.ETypeClick)
		p.Add(btndelete)

		fmt.Println("Panel Id:", p.Id().String())

		panels = append(panels, p)

		parent.Add(p)
	}
}

type myBtnRefresh struct {
	address string
	port    int
	parent  gwu.Panel
	panels  []gwu.Panel
}

func (h *myBtnRefresh) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {

		fmt.Println("Refresh called")
		fmt.Println("address:", h.address)
		fmt.Println("port:", h.port)
		fmt.Println("parent:", h.parent.Id().String())
		for i := 0; i < len(h.panels); i++ {
			fmt.Println("panel:", h.panels[i].Id().String())
		}

		refresh(h.address, h.port, h.parent, h.panels)

		e.MarkDirty(h.parent)
	}
}

func listuserform(address string, port int) gwu.Panel {
	acctlist := gwu.NewVerticalPanel()
	acctlist.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	acctlist.SetCellPadding(10)
	acctlist.SetAttr("width", "800")

	var panels []gwu.Panel

	title2 := gwu.NewHorizontalPanel()
	title2.SetCellPadding(2)
	title2.Add(gwu.NewLabel("Current List of Accounts"))
	btnrefresh := gwu.NewButton("Refresh")
	btnrefresh.SetAttr("align", "center")
	btnrefresh.AddEHandler(&myBtnRefresh{address, port, acctlist, panels}, gwu.ETypeClick)
	title2.Add(btnrefresh)

	acctlist.Add(title2)

	//get current list of accounts
	refresh(address, port, acctlist, panels)

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

	//Add users
	win.Add(adduserform())

	//Display users...
	win.Add(listuserform(address, port))

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("", "127.0.0.1:8000")
	server.SetText("Test GUI App")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser
}
