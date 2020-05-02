package agent

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	port = getPort()
)

func getPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	p := listener.Addr().(*net.TCPAddr).Port
	return fmt.Sprintf(":%v", p)
}

var sampleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello HTTP Test")
})

func TestHttpNormal(t *testing.T) {
	ts := httptest.NewServer(sampleHandler)
	defer ts.Close()

	// リクエストの送信先はテストサーバのURLへ。
	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "Hello HTTP Test" != string(data) {
		t.Fatalf("Data Error. %v", string(data))
	}
}

func TestCountHandler(t *testing.T) {
	// handler := new(countHandler).ServeHTTP
	req := httptest.NewRequest("GET", "/count", nil)
	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(HealthCheckHandler)
	// handler(w, req)
	handler := &countHandler{}
	// new(countHandler).ServeHTTP(w, req)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	t.Log(resp.StatusCode)
	t.Log(resp.Header.Get("Content-Type"))
	t.Log(string(body))
}

func TestUploadHandler(t *testing.T) {
	nodeDir, _ := ioutil.TempDir("", "node")
	defer os.RemoveAll(nodeDir)
	uploadhandler := &uploadHandler{
		NodeDir: nodeDir,
	}
	fieldname := "file"
	filename := "../testdata/ptune/gcagent.toml"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("file not found. %v", err)
	}
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(fieldname, filename)
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatalf("file copy. %v", err)
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		t.Fatalf("file copy. %v", err)
	}
	url := "/upload?node=ostrich&platform=linuxconf"
	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	uploadhandler.ServeHTTP(w, req)
	// handler(w, req)
	resp := w.Result()
	body2, _ := ioutil.ReadAll(resp.Body)
	t.Log(resp.StatusCode)
	t.Log(resp.Header.Get("Content-Type"))
	t.Log(string(body2))
}
