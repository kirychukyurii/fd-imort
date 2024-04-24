package httpserver

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
	"github.com/kirychukyurii/fd-import/models"
	"github.com/kirychukyurii/fd-import/pkg/db"
)

type Attachment struct {
	cfg    *config.Config
	log    *wlog.Logger
	dbpool *db.Connection
}

func NewAttachmentHandler(cfg *config.Config, log *wlog.Logger, dbpool *db.Connection) *Attachment {
	return &Attachment{
		cfg:    cfg,
		log:    log,
		dbpool: dbpool,
	}
}

func (a *Attachment) Attachment(w http.ResponseWriter, req *http.Request) {
	domainID, err := parseInt(req.PathValue("domain_id"))
	if err != nil {
		JSON(w, Error{Msg: "domain id is invalid"}, http.StatusBadRequest)

		return
	}

	_, err = parseInt(req.PathValue("ticket_id"))
	if err != nil {
		JSON(w, Error{Msg: "ticket id is invalid"}, http.StatusBadRequest)

		return
	}

	id, err := parseInt(req.PathValue("id"))
	if err != nil {
		JSON(w, Error{Msg: "attachment id is invalid"}, http.StatusBadRequest)

		return
	}

	attachment, err := a.dbpool.Attachment(req.Context(), domainID, id)
	if err != nil {
		JSON(w, Error{Msg: fmt.Sprintf("database: %s", err)}, http.StatusInternalServerError)

		return
	}

	domain, err := a.dbpool.Domain(req.Context(), &models.Domain{ID: domainID})
	if err != nil {
		JSON(w, Error{Msg: fmt.Sprintf("database: %s", err)}, http.StatusInternalServerError)

		return
	}

	fileExt := filepath.Ext(attachment.Name)
	fileName := strconv.FormatInt(attachment.ID, 10)
	if fileExt != "" {
		fileName = fmt.Sprintf("%d%s", attachment.ID, fileExt)
	}

	file := filepath.Join(a.cfg.AttachmentDir, domain.Name, req.PathValue("ticket_id"), fileName)
	f, err := getFile(file)
	if err != nil {
		JSON(w, Error{Msg: fmt.Sprintf("get file: %s", err)}, http.StatusInternalServerError)

		return
	}

	headers := http.Header{}
	headers.Add("Content-Type", attachment.ContentType)
	headers.Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", encodeURIComponent(attachment.Name)))
	if w.Header().Get("Content-Encoding") == "" {
		headers.Add("Content-Length", strconv.FormatInt(attachment.FileSize, 10))
	}

	File(w, f, http.StatusOK, headers)
}

func getFile(file string) ([]byte, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	return f, nil
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}
