package controllers

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
	"tmp/framework"
)

func StudentsController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	name := ctx.QueryValue("name")

	ctx.WriteJson(struct {
		Name string `json:"name"`
	}{
		Name: name,
	})
}

func ListController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	list := make([]string, 0)

	ctx.WriteString(list[3])
}

func UsersController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	time.Sleep(time.Second * 5)

	ctx.WriteString("users")
}

func ListItemController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	idStr := ctx.PathParam(":item_id")
	name := ctx.PathParam(":item_name")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		ctx.WriteHeader(500)
		return
	}

	ctx.WriteJson(struct {
		Id   int    `json:"item_id"`
		Name string `json:"item_name"`
	}{
		Id:   id,
		Name: name,
	})
}

type PostsPageForm struct {
	Name string
}

func TopPageController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	authUser := ctx.GetValue("AuthUser", 3)

	user := &PostsPageForm{
		Name: authUser.(string),
	}

	ctx.ExecuteTemplate("index.html", user)
}

func PostController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	name := ctx.FormValue("name")
	age := ctx.FormValue("age")
	fileInfo, err := ctx.FormFile("file")
	if err != nil {
		log.Println(err)
		ctx.WriteHeader(500)
	} else {
		err = os.WriteFile(fmt.Sprintf("static/%s_%s_%s", name, age, fileInfo.Filename), fileInfo.Data, fs.ModePerm)
		if err != nil {
			log.Println(err)
			ctx.WriteHeader(500)
		}
	}

	ctx.WriteString("save file")
}

type UserPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UserPostController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	userPost := &UserPost{}
	if err := ctx.BindJson(userPost); err != nil {
		ctx.WriteHeader(500)
	}

	ctx.WriteJson(userPost)
}

func JsonpController(ctx *framework.MyContext) {
	if ctx.Timeout() {
		return
	}

	callback := ctx.QueryValue("callback")
	ctx.JsonP(callback, "Hello, World!")
}
