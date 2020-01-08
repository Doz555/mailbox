package main

//import . "llrp"
//import "os"
//import "io"
import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/mux"
import "github.com/gorilla/sessions"
import "log"
import "net/http"

import "time"

var store = sessions.NewCookieStore([]byte("nirvana")) //นิพพาน
var constatus = false
var connum = 2
var con = "root:@tcp(localhost:3306)/warehouse_binary2"

type sc struct {
	br string
}

func loginhtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/login.html")
}
func registerhtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/register.html")
}
func login(w http.ResponseWriter, r *http.Request) {
	
	psdd := r.FormValue("psw")
	unamee := r.FormValue("uname")
	if unamee != "" && psdd != "" {
		qr := `select COUNT(*) from user Where UserName = "` + unamee + `" and Pass = PASSWORD("` + psdd + `")`
		json := selectvalue(qr)
		if json == "1" {
			//nameshow:=selectvalue(`select Name from user Where UserName = "` + unamee + `"`)
			sessionSD(w, r, "login", unamee)
			//fmt.Fprintln(w,`<html><style></style><body><div align ="right"><a name="uname"> `+ nameshow +`</a> <a href="/src/login.html"><Button>logout </Button></a><a href="/src/register.html"><Button>Register </Button></a>	</div></body></html>`)
			fmt.Fprintln(w,`<script>window.location.assign("src/Main.html")</script>`)
		}else {
			sessionSD(w,r,"login","NULL") 
			fmt.Fprintln(w, `<script>window.location.assign(\"login\")</script>`)
		}
	}
}
func Register(w http.ResponseWriter, r *http.Request) {
	RGname := r.FormValue("RGname")
	RGtell := r.FormValue("RGtell")
	RGuname := r.FormValue("RGuname")
	RGpsw := r.FormValue("RGpsw")
	sessionSD(w,r,"login","NULL") 
	if RGname != "" && RGtell != "" && RGuname != "" && RGpsw != "" {
		NameCK := selectvalue(`select COUNT(*) from user Where UserName = "` + RGname +`"`)
		UNameCK := selectvalue(`select COUNT(*) from user Where UserName = "` + RGuname +`"`)
		if NameCK == "0" && UNameCK == "0" {
			selectvalue(`INSERT INTO user (ID_User, Name, Tel, UserName, Pass) VALUES (NULL, "` + RGname + `", "` + RGtell + `", "` + RGuname + `", PASSWORD("` + RGpsw + `"));`)

			fmt.Fprintln(w, `<script>window.location.assign("src/login.html")</script>`)
		} else {
			fmt.Fprintln(w, `<script>window.location.assign("src/register.html")</script>`)
			
		}
	}
}
func product(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, selectjson(`SELECT concat('{',group_concat('"',ID_Pro,'":{"Serial":"',Serial,'","ID_Model":"',ID_Model,'"}'),'}') as json FROM product`))
	}						
func logout(w http.ResponseWriter, r *http.Request) {
	sessionSD(w,r,"login","NULL") 
	fmt.Fprintln(w,"<script>window.location.assign(\"/src/login.html\")</script>")
}

func sessionSD(w http.ResponseWriter, r *http.Request, name string, data string) {

	session, err := store.Get(r, "yms")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values[name] = data
	session.Save(r, w)
}
func sessionGD(w http.ResponseWriter, r *http.Request, name string) string {

	session, err := store.Get(r, "yms")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "session eror"
	}
	data := ""
	if session.Values[name] == nil {
		data = "NULL"
	} else {
		data = session.Values[name].(string)
	}
	session.Save(r, w)
	return data
}
func selectvalue(qr string) string {

	db, err := sql.Open("mysql", con)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	//print(qr)
	rows, err := db.Query(qr)
	if err != nil {
		fmt.Printf("park")
		panic(err.Error())
	}
	defer rows.Close()
	value := ""
	rows.Next()
	rows.Scan(&value)
	return value

}
func selectjson(qr string) string {
	
	db, err := sql.Open("mysql", con)
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	
	rows, err := db.Query(qr)
	if err != nil {	fmt.Printf("park") ; panic(err.Error())}
	defer rows.Close()
	json:="";rows.Next();rows.Scan(&json)
	return json
	
}
func checklogin(w http.ResponseWriter, r *http.Request)bool{
	if( sessionGD(w,r,"login") == "NULL"){
		fmt.Fprintln(w,"<script>window.location.assign(\"/src/login.html\")</script>")
		return false	
	}else {return true }
}	
func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("%s", time.Now().Sub(start).String)
	}()

	r := mux.NewRouter()
	r.HandleFunc("/Register.php", Register)
	r.HandleFunc("/login", loginhtml)
	r.HandleFunc("/product.php", product)
	r.HandleFunc("/register", registerhtml)
	r.HandleFunc("/login.php", login)
	r.PathPrefix("/src/").Handler(http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))
	r.PathPrefix("/pic/").Handler(http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic/"))))
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("./bower_components/"))))

	// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Printf("start server 80 ")
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
