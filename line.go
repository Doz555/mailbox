package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"


	"github.com/gorilla/mux"
	"github.com/utahta/go-linenotify"
)

type MailBox struct {
	ID   int
	ID2  int
	Mode string
	Data string
}

var passWord [4]string
var con [4]net.Conn
var constatus = false
var connum = 2
var token = "Q9DFsGbBMhlTB2zWRWGecCQR9MTfxE3QYBngpMasMki"


func tcpSrv() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}
		go handleConnection(conn)
	}
}

type AppPi struct {
	IP   string
	Pass string
	User string
}

func OpenDoor(w http.ResponseWriter, r *http.Request) {
	var appPi AppPi
	json.NewDecoder(r.Body).Decode(&appPi)
	i, err := strconv.Atoi(appPi.IP)
	if err != nil {
		log.Println(err)
	}
	if con[i] == nil {
		fmt.Fprintln(w, `{"steus":"noCon"}`)
		return
	}
	if appPi.User != "fail" {
		if passWord[i] == appPi.Pass {
			con[i].Write([]byte("pass" + "\r\n"))
			fmt.Fprintln(w, `{"steus":"OK"}`)
		} else {
			con[i].Write([]byte("worngpass" + "\r\n"))
			fmt.Fprintln(w, `{"steus":"fail"}`)
		}
	} else {
		fmt.Fprintln(w, `{"steus":"login"}`)
	}
}
func notification(w http.ResponseWriter, r *http.Request) {
	MSG := r.FormValue("msg")
	c := linenotify.New()
	c.Notify(token, MSG, "", "", nil)
	fmt.Fprintln(w, `ok`)
}
func notification2(MSG string) {
	c := linenotify.New()
	c.Notify(token, MSG, "", "", nil)
}
func node(IP int, MSG string) {
	con[IP].Write([]byte(MSG + "\r\n"))
}

func main() {
	go tcpSrv()
	print(".......< Loading >.......")
	start := time.Now()
	defer func() {
		fmt.Printf("%s", time.Now().Sub(start).String)
	}()
	r := mux.NewRouter()
	r.HandleFunc("/OpenDoor.php", OpenDoor)
	r.HandleFunc("/getBoxST.php", getBoxST)
	r.HandleFunc("/getMailST.php", getMailST)
	r.HandleFunc("/login.php", LoginJson)
	r.HandleFunc("/notification", notification)
	fmt.Printf("start server")
	println("")
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
func handleConnection(conn net.Conn) {

	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("emty data")
		return
	}
	message := string(bufferBytes)
	clientAddr := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + " from " + clientAddr)
	var mailBox MailBox
	json.NewDecoder(strings.NewReader(message)).Decode(&mailBox)
	con[mailBox.ID] = conn
	if mailBox.Mode == "Connect" {
		notification2("MailBox Connect")
	} else if mailBox.Mode == "Line" {
		notification2(mailBox.Data)
	} else if mailBox.Mode == "pass" {
		notification2("Password :: " + mailBox.Data)
		passWord[mailBox.ID] = mailBox.Data
	} else if mailBox.Mode == "node" {
		node(mailBox.ID2, mailBox.Data)
		if "openDoor" == mailBox.Data {
			s := strconv.Itoa(mailBox.ID)
			selectvalue(`INSERT INTO boxST (id_stetus, id_mailbox, date) VALUES (NULL, '` + s + `', current_timestamp());`)
			selectvalue(`UPDATE mailST SET stetus = 1 WHERE id_mailbox = ` + s)
		}
	} else if mailBox.Mode == "NewMSG" {
		notification2("You have new Mail")
		node(mailBox.ID2, mailBox.Data)
		s := strconv.Itoa(mailBox.ID2)
		selectvalue(`INSERT INTO mailST (id_stetus, id_mailbox, stetus, date) VALUES (NULL, '` + s + `', '0', current_timestamp());`)
	}
	log.Println(response)
	conn.Write([]byte("Connect Success" + "\r\n"))
	//handleConnection(conn)
}
func GetObjJsonFromRequest(w http.ResponseWriter, r *http.Request, result *map[string]interface{}) {
	json.NewDecoder(r.Body).Decode(&result)
}
