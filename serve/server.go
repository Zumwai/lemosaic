package serve

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	//"image"
	//"mime/multipart"
	"image/png"
	"mosaic/config"
	"mosaic/imgConv"
	"mosaic/localMosaic"
	"net/http"
	"time"
	//"strconv"
)

/*running server, duh */
func StartServer() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", upload)
	mux.HandleFunc("/mosaic", mosaic)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	//TILESDB = tilesDB()
	fmt.Println("Mosaic server started.")
	server.ListenAndServe()
}

/* loads html page with load image button */
func upload(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.html")
	t.Execute(w, nil)
}

/*main html mosaic. Doesn't uses regular tools, so it's an ugly bastard*/
func mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	r.ParseMultipartForm(1e+7)
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	//var te multipart.File
	//tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	//original, _, _ := image.Decode(file)
	or, err := localMosaic.DecodeByType("image/png", file)
	if err != nil {
		fmt.Println(err)
		return
	}
	original, err := imgConv.ResizeInMemory(or, or.Bounds().Max.X, or.Bounds().Max.Y)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	hashed, err := localMosaic.PopulateHashDir("./pics/")
	if err != nil {
		fmt.Println(err)
		return
	}
	final := imgConv.PrepareMosaic(original, config.ChunkLookup(), config.RoutineLookup(), hashed)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, final)
	if err != nil {
		fmt.Println(err)
		return
	}
	mos := base64.StdEncoding.EncodeToString(buf.Bytes())
	//fmt.Println(mos)
	t1 := time.Now()
	images := struct {
		mosaic   string
		duration string
	}{mos, fmt.Sprintf("%v ", t1.Sub(t0))}

	t, err := template.ParseFiles("result.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, images)
}
