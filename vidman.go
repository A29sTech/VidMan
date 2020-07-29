package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/A29sTech/VidMan/core"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // This Is Ther Driver Of SQLTILE3 Only
)

// Some Importent Constents.
const (
	VideoFormFeild = "video"
	DocFormFeild   = "doc"
)

// Globaly Acsessable Vars ;
var (
	lib string
	db  *gorm.DB
)

///////// MOCK DATA ////////
// vid := core.ViData{
// 	Title:    "Hello World",
// 	Subject:  "he",
// 	Desc:     "am",
// 	Tags:     "niodec",
// 	Filename: "nonamed",
// 	Docname:  "nodeoc",
// }

func main() {

	var port string // Http server port.
	var home string
	var isAdmin bool

	// Parse Command Line.
	flag.StringVar(&lib, "lib", "", "-lib 'videopath' including *.db")
	flag.StringVar(&port, "port", "3333", "-port 'http-server-port'")
	flag.StringVar(&home, "home", "", "-home 'index.html file require in home dir.")
	flag.BoolVar(&isAdmin, "admin", false, "-admin , for write, update and delete Acsess.")
	flag.Parse()

	// Lib Arg Check ;
	if lib == "" {
		fmt.Println("lib path is not provided as, -lib libpath")
		return
	}

	// Setup DataBase ;
	var err error
	if db, err = core.OpenViDB("Videos.db"); err != nil {
		fmt.Println(err)
		return
	}

	// Create Mux Router ;
	router := mux.NewRouter()

	// Register Video Files Server ;
	router.PathPrefix("/cdn/").Handler(http.StripPrefix("/cdn/",
		http.FileServer(http.Dir(lib))))

	// Register Api for All.
	router.HandleFunc("/api/Video/{id:[0-9]+}", getVideoAPI)
	router.HandleFunc("/api/SearchVideos/{colum}/{query}/{limit:[0-9]+}/{offset:[0-9]+}", searchVideoAPI)
	router.
		HandleFunc("/api/Videos/{limit:[0-9]+}/{offset:[0-9]+}", allVideosAPI).
		Methods(http.MethodGet)

	// Register if Admin Permission Given ;
	if isAdmin {
		router.HandleFunc("/api/AddVideo", addVideoAPI).Methods(http.MethodPost)
		router.HandleFunc("/api/UpdateVideo/{id:[0-9]+}", updateVideoAPI).Methods(http.MethodPost)
		router.HandleFunc("/api/DeleteVideo/{id:[0-9]+}", deleteVideoAPI)
	}

	// Get By Tags, Like PlayList Altranative;
	router.HandleFunc("/api/PlayListTags/{tag}", getByTegs)

	// Check For Home and Serve Html UI;
	if home != "" {
		router.HandleFunc("/", homeHandler(path.Join(home, "index.html")))
		router.PathPrefix("/").
			Handler(FileServerWithCustom404(http.Dir(home),
				path.Join(home, "index.html")))
		// router.NotFoundHandler = http.HandlerFunc(homeHandler(path.Join(home, "index.html")))
	}

	// Running Message
	fmt.Println("Server Started : http://127.0.0.1:" + port)

	// Start The Server.
	http.ListenAndServe("127.0.0.1:"+port, router)

}

// FileServerWithCustom404 : For SPA ;
func FileServerWithCustom404(fs http.FileSystem, entrypoint string) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			// Call Not Found Handler ;
			homeHandler(entrypoint)(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

// IndexHandler For Static Html, Javascript, Css Fire Server
func homeHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return http.HandlerFunc(fn)
}

// Get Video By Id ;
func getVideoAPI(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	vidata := core.ViData{}
	if err := db.Table(core.TableName).First(&vidata, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
	}
	jsonBytes, err := json.Marshal(vidata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes))

}

// Search Videos By Colum Name ; Default is Title;
func searchVideoAPI(w http.ResponseWriter, r *http.Request) {

	// Retrive 'colum', 'query', 'limit', 'offset' from url param ;
	params := mux.Vars(r)
	columName, query := params["colum"], params["query"]
	limit, offset := params["limit"], params["offset"]

	// Craete a ViData Slice ;
	visdata := []core.ViData{}

	// Check geven colum name exist in database ;
	if !core.ViDataHasColumName(columName) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ViData Table Does Not Contain Colum Name : %s", columName)
		return
	}

	// Query Matched Data form DataBase to 'vidata' Slice ;
	err := db.Table(core.TableName).
		Where(columName+" LIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Order(columName).
		Find(&visdata).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode Result To Json ;
	jsonBytes, err := json.Marshal(visdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send Json Result ;
	fmt.Fprintf(w, "%s", string(jsonBytes))

}

// allVideosAPI : get all videos by limit & Offset ;
func allVideosAPI(w http.ResponseWriter, r *http.Request) {

	// Retrive 'limit' & 'offset' from url param ;
	params := mux.Vars(r)
	limit, offset := params["limit"], params["offset"]

	// Create a ViData Slice ;
	vidataSlice := []core.ViData{}

	// Retrive All Matched Query To 'vidataSlice' ;
	err := db.Table(core.TableName).Offset(offset).Limit(limit).Find(&vidataSlice).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Encode 'vidata' to Json ;
	jsonBytes, err := json.Marshal(&vidataSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Finally Send The Json Results ;
	fmt.Fprintf(w, string(jsonBytes))

}

// addVideoApi : Add A Video To Lib.
func addVideoAPI(w http.ResponseWriter, r *http.Request) {

	// Parse Multipart Form Data.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Make A VideoModel From Request Form ;
	vidata, err := core.Form2ViData(r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Video File Handling ...
	file, header, err := r.FormFile(VideoFormFeild)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Save File And Get FileName ;
	filename := header.Filename
	err = core.SaveUploadedFile(file, path.Join(lib, filename))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Finally Add Video To Lib ;
	vidata.Filename = filename
	if err = db.Table(core.TableName).Create(&vidata).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		os.Remove(path.Join(lib, filename))
		return
	}

	// Encode ViData to Json ;
	jsonBytes, err := json.Marshal(vidata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes))
}

// update Databse; Require Param : id & json body ;
func updateVideoAPI(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	// Convert ID string to int ;
	intID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Create A ViData Struct, to store json req data;
	vidata := core.ViData{}
	vidata.ID = intID

	// Pasre Body To Json ;
	err = json.NewDecoder(r.Body).Decode(&vidata)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Update To DB ;
	err = db.Table(core.TableName).Model(&vidata).Omit("filename").Update(&vidata).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Retrive Updated Data From Database ;
	err = db.Table(core.TableName).First(&vidata, intID).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Encode ViData to Json ;
	jsonBytes, err := json.Marshal(vidata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Send Encoded Json Data ;
	fmt.Fprintf(w, "%s", string(jsonBytes))

}

func getByTegs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tag := params["tag"]

	// Craete a ViData Slice ;
	visdata := []core.ViData{}

	// Query Matched Data form DataBase to 'vidata' Slice ;
	err := db.Table(core.TableName).
		Where("tags LIKE ?", "%"+tag+"%").
		Order("indx").
		Find(&visdata).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode Result To Json ;
	jsonBytes, err := json.Marshal(visdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send Json Result ;
	fmt.Fprintf(w, "%s", string(jsonBytes))
}

// deleteVideoAPI : To Delete A Video ;
func deleteVideoAPI(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	// Convert ID string to int ;
	intID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	// Delete Op ;
	err = db.Table(core.TableName).Delete(core.ViData{ID: intID}).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	fmt.Fprintf(w, "%s", id)

}
