// API is just a small program to host the demo API calls
package main

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// services
	"github.com/metarex-media/mrx-demo-svc/api/battery"
	"github.com/metarex-media/mrx-demo-svc/api/gps"
	"github.com/metarex-media/mrx-demo-svc/api/mrx"
	"github.com/metarex-media/mrx-demo-svc/api/ninjs"
	"github.com/metarex-media/mrx-demo-svc/api/qc"
	"github.com/metarex-media/mrx-demo-svc/api/rnf"
	"github.com/metarex-media/mrx-demo-svc/api/wavdraw"
)

// ErrorMessage is the error message format returned by the API as json
type ErrorMessage struct {
	Error string `json:"error"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Battery Paths
	e.POST("/battery", HandlerBuild(battery.ToPNG, "image/png"))
	//	e.POST("/batteryStagger", HandlerBuild(battery.BatteryToPNGStagger, "image/png"))
	e.POST("/batteryFault", HandlerBuild(battery.FaultToJPEG, "image/jpeg"))

	// NinJS Paths
	e.POST("/ninjsToMD", HandlerBuild(ninjs.ToMD, echo.MIMEApplicationJSON))
	e.POST("/ninjsToNewsml", HandlerBuild(ninjs.ToNewsMl, echo.MIMEApplicationXML))

	// QC
	e.POST("/qcToGraph", HandlerBuild(qc.GenBarChart, "image/png"))

	// MXF Extract
	//	e.POST("/mxfContents", HandlerBuild(mrx.MXFHeaderContents, echo.MIMEApplicationJSON))
	e.POST("/C2PAExtract", HandlerBuild(mrx.JpegC2PA, echo.MIMEApplicationJSON))

	// Dawn Chorus
	e.POST("/gps", HandlerBuild(gps.ConvertGPX, echo.MIMEApplicationJSON))
	e.POST("/waveform", HandlerBuild(wavdraw.Visualise, "image/png"))

	// RNF demos
	e.POST("/ffmpeg", HandlerBuild(rnf.GenerateFFmpeg, echo.MIMETextPlain, rnf.GetFFmpegParams()...))

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

// APIHandle is format for all the api functions to be built as part of the API.
type APIHandle func(requestBody []byte, apiParams ...string) ([]byte, error)

// HandlerBuild generates a boilerplate Echo Handler
// the parameter titles are in the order they are expected
// in the sub function
func HandlerBuild(apiFunc APIHandle, outputMime string, params ...string) echo.HandlerFunc {
	return func(c echo.Context) error {

		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
			if len(bodyBytes) == 0 {
				return c.JSON(http.StatusBadRequest, ErrorMessage{"no data received"})
			}
		} else {
			return c.JSON(http.StatusBadRequest, ErrorMessage{"no data received"})
		}

		// get all required parameters
		apiParams := make([]string, len(params))
		for i, param := range params {
			apiParams[i] = c.QueryParam(param)
		}

		output, err := apiFunc(bodyBytes, apiParams...)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorMessage{err.Error()})
		}

		return c.Blob(http.StatusOK, outputMime, output)
	}
}
