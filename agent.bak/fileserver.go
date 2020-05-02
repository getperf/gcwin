package agent

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/getperf/gcagent/https"
	log "github.com/sirupsen/logrus"
)

const defaultUrl = "0.0.0.0:59000"

type FileServer struct {
	Address     string `mapstructure:"server_url"`
	UrlPrefix   string `mapstructure:"url_prefix"`
	DownloadDir string `mapstructure:"serve_dir"`
	NodeDir     string `mapstructure:"node_dir"`
	TlsConfig   string `mapstructure:"tls_config"`
}

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

type uploadHandler struct {
	mu      sync.Mutex
	NodeDir string
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

// func hundleUpload(w http.ResponseWriter, r *http.Request) {
func (h *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	node, ok := r.Form["node"]
	if !ok {
		http.Error(w, "node parameter not found", 500)
		return
	}
	platform, ok := r.Form["platform"]
	if !ok {
		http.Error(w, "platform parameter not found", 500)
		return
	}
	log.Info("node :", node[0], ",platform : ", platform[0])
	h.mu.Lock()
	defer h.mu.Unlock()
	// "file"というフィールド名に一致する最初のファイルが返却される
	// マルチパートフォームのデータはパースされていない場合ここでパースされる
	formFile, _, err := r.FormFile("file")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	defer formFile.Close()

	// データを保存するファイルを開く
	nodedir := filepath.Join(h.NodeDir, node[0])
	if err := os.MkdirAll(nodedir, 0777); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	nodefile := fmt.Sprintf("%s.json", platform[0])
	nodepath := filepath.Join(nodedir, nodefile)
	saveFile, err := os.Create(nodepath)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	defer saveFile.Close()

	// ファイルにデータを書き込む
	_, err = io.Copy(saveFile, formFile)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Info("save ", nodepath)

	w.WriteHeader(http.StatusCreated)
}

func (fs *FileServer) serve() {

	service := http.StripPrefix(fs.UrlPrefix,
		http.FileServer(http.Dir(fs.DownloadDir)))

	http.Handle(fs.UrlPrefix, service)
	http.Handle("/count", new(countHandler))
	handleUpload := &uploadHandler{
		NodeDir: fs.NodeDir,
	}
	http.Handle("/upload", handleUpload)
	// http.HandleFunc("/upload", handleUpload)

	server := &http.Server{
		Addr: fs.Address,
	}
	https.Listen(server, fs.TlsConfig)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (config *Config) FileServe() {
	serviceAddress := defaultUrl
	schedule := config.Schedule
	if schedule != nil && schedule.ServiceUrl != "" {
		serviceAddress = schedule.ServiceUrl
	}
	fs := &FileServer{
		Address:     serviceAddress,
		UrlPrefix:   "/config/",
		DownloadDir: config.ArchiveDir,
		NodeDir:     config.NodeDir,
	}
	log.Info("set download service : ", fs)
	go fs.serve()
}
