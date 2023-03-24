package serve

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"mosaic/config"
	"mosaic/imgConv"
	"mosaic/localMosaic"
	"net/http"
	"time"
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
	fmt.Println("Mosaic server started.")
	server.ListenAndServe()
}

/* loads html page with load image button */
func upload(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("upload.html")
	if err != nil {
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		return
	}
}

/*main html mosaic. Doesn't uses regular tools, so it's an ugly bastard*/
func mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	r.ParseMultipartForm(1e+7)
	file, _, err := r.FormFile("image")
	config.SetChunkSize(r.FormValue("tile-size"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	unformatted, err := localMosaic.DecodeByType(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	src := imgConv.ConvertToDrawable(unformatted)
	if err != nil {
		fmt.Println(err)
		return
	}

	hashed, err := localMosaic.PopulateHashDir(config.SrcImagesLookup())
	if err != nil {
		fmt.Println(err)
		return
	}
	final := imgConv.PrepareMosaic(src, hashed)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, final)
	if err != nil {
		fmt.Println(err)
		return
	}

	mos := base64.StdEncoding.EncodeToString(buf.Bytes())
	t1 := time.Now()
	images := struct {
		Mosaic   string
		Duration string
	}{mos, fmt.Sprintf("%v ", t1.Sub(t0))}

	t, err := template.ParseFiles("result.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(w, images)
	if err != nil {
		fmt.Println(err)
	}
}
