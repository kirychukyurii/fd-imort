package httpserver

import (
	"net/http"

	"github.com/kirychukyurii/fd-import/pkg/db"
)

func (s *Server) RegisterHandlers(dbpool *db.Connection) {
	attachment := NewAttachmentHandler(s.cfg, s.log, dbpool)
	s.router.HandleFunc("GET /ping", func(w http.ResponseWriter, req *http.Request) {
		JSON(w, "ok", http.StatusOK)
	})

	s.router.HandleFunc("GET /{domain_id}/ticket/{ticket_id}/attachments/{id}", attachment.Attachment)
}
