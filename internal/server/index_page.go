package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/tasansga/grantory/internal/storage"
)

var (
	//go:embed templates/index.html
	indexTemplateSource string
	//go:embed static/water.min.css
	waterCSS []byte

	indexTemplate = template.Must(template.New("index").Funcs(template.FuncMap{
		"labelSummary": labelSummary,
	}).Parse(indexTemplateSource))
)

type indexPageData struct {
	Namespace            string
	RequestsWithGrant    int64
	RequestsWithoutGrant int64
	TotalRequests        int64
	TotalGrants          int64
	TotalRegisters       int64
	Hosts                []storage.Host
	Requests             []storage.Request
	Grants               []storage.Grant
	Registers            []storage.Register
}

func (s *Server) handleIndex(c *fiber.Ctx) error {
	logRequestEntry(c, "Server.handleIndex", nil)

	store, namespace, err := resolveNamespaceStore(c)
	if err != nil {
		return err
	}

	reqCounts, err := store.CountRequestsByGrantPresence(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("count requests for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to collect request stats")
	}

	registerCounts, err := store.CountRegisters(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("count registers for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to collect register stats")
	}

	grantCounts, err := store.CountGrants(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("count grants for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to collect grant stats")
	}

	hosts, err := store.ListHosts(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("list hosts for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to list hosts")
	}

	requests, err := store.ListRequests(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("list requests for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to list requests")
	}

	registers, err := store.ListRegisters(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("list registers for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to list registers")
	}

	grants, err := store.ListGrants(c.Context())
	if err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("list grants for index")
		return fiber.NewError(http.StatusInternalServerError, "unable to list grants")
	}

	data := indexPageData{
		Namespace:            namespace,
		RequestsWithGrant:    reqCounts["with_grant"],
		RequestsWithoutGrant: reqCounts["without_grant"],
		TotalRequests:        reqCounts["with_grant"] + reqCounts["without_grant"],
		TotalGrants:          grantCounts["total"],
		TotalRegisters:       registerCounts["total"],
		Hosts:                hosts,
		Requests:             requests,
		Grants:               grants,
		Registers:            registers,
	}

	var buf bytes.Buffer
	if err := indexTemplate.Execute(&buf, data); err != nil {
		logrus.WithError(err).WithField("namespace", namespace).Error("render index page")
		return fiber.NewError(http.StatusInternalServerError, "unable to render stats page")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Status(http.StatusOK).Send(buf.Bytes())
}

func (s *Server) handleWaterCSS(c *fiber.Ctx) error {
	logRequestEntry(c, "Server.handleWaterCSS", nil)
	c.Set(fiber.HeaderContentType, "text/css; charset=utf-8")
	c.Set(fiber.HeaderCacheControl, "public, max-age=31536000")
	return c.Status(http.StatusOK).Send(waterCSS)
}

func labelSummary(labels map[string]string) string {
	if len(labels) == 0 {
		return "—"
	}
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, labels[key]))
	}
	return strings.Join(parts, ", ")
}
