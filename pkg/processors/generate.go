package processors

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
	"gitlab.com/joshraphael/motdoftheday/settings"
)

type generatedPost struct {
	Post       *database.Post
	User       *database.User
	LatestPost *database.PostHistory
	Categories []database.Category
	Tags       []database.Tag
}

func (prcr Processor) generatePost(p post.Post) apierror.IApiError {
	db_post, err := prcr.db.GetPostByUrlTitle(p.UrlTitle())
	if err != nil {
		msg := "error getting post " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	if db_post == nil {
		msg := "no post found " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	user, err := prcr.db.GetUserById(db_post.UserID)
	if err != nil {
		msg := "error getting userwhen generating post: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	if user == nil {
		fmt.Println("tester")
		msg := "no user found when generating post"
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	latest_post, err := prcr.db.GetLatestPost(db_post)
	if err != nil {
		msg := "error getting latest post " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	if latest_post == nil {
		msg := "no post history found " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	categories, err := prcr.db.GetPostCategories(latest_post)
	if err != nil {
		msg := "error getting post categories " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	if len(categories) == 0 {
		msg := "no categories for post " + p.UrlTitle()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	tags, err := prcr.db.GetPostTags(latest_post)
	if err != nil {
		msg := "error getting post tags " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	if len(tags) == 0 {
		msg := "no tags for post " + p.UrlTitle()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	if db_post == nil {
		msg := "Post does not exists " + p.UrlTitle() + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", p.Method())
		return apiErr
	}
	if _, err := os.Stat(settings.POST_TEMPLATE); err != nil {
		msg := "Template file " + settings.POST_TEMPLATE + " does not exist: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	tmpl, err := template.ParseFiles(settings.POST_TEMPLATE)
	if err != nil {
		msg := "Cannot read template file " + settings.POST_TEMPLATE + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	year, month, day := time.Unix(latest_post.InsertTime, 0).UTC().Date()
	if _, err := os.Stat(settings.POST_DIR); os.IsNotExist(err) {
		e := os.MkdirAll(settings.POST_DIR, os.ModePerm)
		if e != nil {
			msg := "cannot create post dir " + settings.POST_DIR + ": " + e.Error()
			apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
			return apiErr
		}
	}
	filename := settings.POST_DIR + "/" + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + "-" + db_post.UrlTitle + ".md"
	gp := generatedPost{
		Post:       db_post,
		User:       user,
		LatestPost: latest_post,
		Categories: categories,
		Tags:       tags,
	}
	f, err := os.Create(filename)
	if err != nil {
		msg := "cannot open post file" + filename + ": " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	defer f.Close()
	err = tmpl.Execute(f, gp)
	if err != nil {
		msg := "Cannot render template: " + err.Error()
		apiErr := apierror.New(errors.New(msg), "INTERNAL", p.Method())
		return apiErr
	}
	return nil
}
