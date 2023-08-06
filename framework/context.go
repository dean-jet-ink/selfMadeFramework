package framework

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"net/textproto"
	"sync"
	"text/template"
)

type MyContext struct {
	w          http.ResponseWriter
	r          *http.Request
	pathParams map[string]string
	values     map[string]any
	mux        sync.RWMutex
	hasTimeout bool
	handlers   []func(*MyContext)
	index      int
}

func NewMyContext(w http.ResponseWriter, r *http.Request) *MyContext {
	return &MyContext{
		w:      w,
		r:      r,
		values: make(map[string]any),
		index:  -1,
	}
}

func (ctx *MyContext) Timeout() bool {
	return ctx.hasTimeout
}

func (ctx *MyContext) SetTimeout(hasTimeout bool) {
	ctx.hasTimeout = hasTimeout
}

func (ctx *MyContext) GetValue(key string, defaultValue any) any {
	ctx.mux.RLock()
	defer ctx.mux.RUnlock()

	if value, ok := ctx.values[key]; ok {
		return value
	}

	return defaultValue
}

func (ctx *MyContext) SetValue(key string, value any) {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()

	ctx.values[key] = value
}

func (ctx *MyContext) SetHandlers(handlers []func(*MyContext)) {
	ctx.handlers = handlers
}

func (ctx *MyContext) Next() {
	ctx.index++
	for ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *MyContext) Abort() {
	ctx.index = math.MaxInt8
}

func (ctx *MyContext) WriteJson(data any) {
	m, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		ctx.w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.w.Header().Add("Content-Type", "application/json")
	ctx.w.Write(m)
}

func (ctx *MyContext) BindJson(data any) error {
	jsonData, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		return err
	}

	ctx.r.Body = io.NopCloser(bytes.NewBuffer(jsonData))

	return json.Unmarshal(jsonData, data)
}

func (ctx *MyContext) WriteString(str string) {
	io.WriteString(ctx.w, str)
}

func (ctx *MyContext) Queries() map[string][]string {
	return ctx.r.URL.Query()
}

func (ctx *MyContext) QueryValue(key string) string {
	return ctx.r.URL.Query().Get(key)
}

func (ctx *MyContext) SetPathParams(pathParams map[string]string) {
	ctx.pathParams = pathParams
}

func (ctx *MyContext) PathParam(key string) string {
	return ctx.pathParams[key]
}

func (ctx *MyContext) WriteHeader(statusCode int) {
	ctx.w.WriteHeader(statusCode)
}

func (ctx *MyContext) ExecuteTemplate(filePath string, data any) {
	templ, err := template.ParseFiles(filePath)
	if err != nil {
		ctx.w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Execute(ctx.w, data)
}

func (ctx *MyContext) FormValue(key string) string {
	return ctx.r.FormValue(key)
}

type FileInfo struct {
	Data     []byte
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
}

func (ctx *MyContext) FormFile(key string) (*FileInfo, error) {
	fileInfo, fileHeader, err := ctx.r.FormFile(key)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(fileInfo)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Data:     data,
		Filename: fileHeader.Filename,
		Header:   fileHeader.Header,
		Size:     fileHeader.Size,
	}, nil
}

func (ctx *MyContext) JsonP(callback string, param any) error {
	ctx.w.Header().Add("Content-Type", "application/javascript")

	callback = template.JSEscapeString(callback)

	_, err := ctx.w.Write([]byte(callback))
	if err != nil {
		return err
	}

	ctx.w.Write([]byte("("))

	m, err := json.Marshal(param)
	if err != nil {
		return err
	}

	_, err = ctx.w.Write(m)
	if err != nil {
		return err
	}

	ctx.w.Write([]byte(")"))

	return nil
}
