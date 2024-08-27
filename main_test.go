package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAutoEltHandler(t *testing.T) {

	_, err := http.Get("http://localhost:9000/")

	// check the services can run
	if err != nil {
		panic(fmt.Sprintf("error connecting to services server at http://localhost:9000/: %v", err))
	}

	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
	slog.SetDefault(slog.New(h))

	e := echo.New()
	paths := []string{"/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fill"}
	body := `{"percentage":80}`

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := autoELT(c)

		Convey("Checking the autoelt post requests", t, func() {
			Convey(fmt.Sprintf("Making a call to %s with a body of %s", path, body), func() {
				Convey("The a 200 code is returned without an error", func() {
					So(err, ShouldBeNil)
					So(rec.Code, ShouldResemble, http.StatusOK)
				})
			})
		})

	}

}
