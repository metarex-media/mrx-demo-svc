package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/mm-eng/mrx-api-demo/mrxhandle/mrxlog"
	"gitlab.com/mm-eng/mrx-api-demo/mrxhandle/mrxlog/utils"
	"gitlab.com/mm-eng/gl-mrx-demo-svc/transformations"
)

func init() {
	// set up logging
	// utils.JSON(&slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}, "./tmp/")
	utils.ColourConsole(&slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})

	// set the flags
	//	transformFlags()
	//	rootCmd.AddCommand(TransformCmd)

}

// ErrorMessage is the same format as in
// gitlab.com/mm-eng/gl-mrx-demo-svc/api/
type ErrorMessage struct {
	Error string `json:"error"`
}

func main() {

	// Echo instance
	e := echo.New()

	// Middleware

	/*	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, http.MethodGet, http.MethodPost},
	}))*/
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORS())

	// services
	e.POST("/autoelt", autoELT)
	e.GET("/getservices", getServices)

	// check the test is running as intended
	e.GET("/test", test)

	// serve schemas here
	e.Static("/schema", "demodata/schema")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

	/*		data := make([]sphericalCoords, 50)

			for i, dat := range data {

				dat.R = rand.Float64() * 100
				dat.Elevation = rand.Float64() * 90
				dat.Azimuth = rand.Float64() * 180
				data[i] = dat
			}
	*/
	/*
			f, _ := os.Create("./demodata/magicDemo.mrx")
			in := make(chan []byte, 1) //len(data))

			//var inputData []any
			//b, _ := os.ReadFile("./demodata/transformedDataVelocity.mrx.json")
			//err = json.Unmarshal(b, &inputData)

			//for _, dat := range data {
			d := `{
		    "firstName": "john",
		    "lastName": "smith",
		    "primary job": "video systems architect",
		    "secondary job": "",
		    "year":2024,
		    "car" : {
		            "make": "lambo",
		            "year made":2020
		    }
		  }`
			//d, _ := yaml.Marshal(dat)
			in <- []byte(d)
			//	}

			go func() {
				for {

					if len(in) == 0 {
						close(in)
						break

					}
				}
			}()

			demoConfig := encode.Configuration{Version: "pre alpha",
				Default:          encode.StreamProperties{StreamType: "Example XYZ data", FrameRate: "1/1", NameSpace: "https://metarex.media/reg/" + "MRX.123.456.789.ads"},
				StreamProperties: map[int]encode.StreamProperties{0: {NameSpace: "MRX.123.456.789.ads"}},
			}
			err := encode.EncodeSingleDataStream(f, in, demoConfig)
			fmt.Println(err)*/

	//rootCmd.SetUsageTemplate("empty" + rootCmd.UsageTemplate())

	// rootCmd.DebugFlags()
	//cobra.CheckErr(rootCmd.Execute())
	// inOut(dataSender{}, "./demodata/xyzDemo.mrx", "MRX.123.456.789.mph")
	// inOut(dataSender{}, "./demodata/xyzDemo.mrx", "MRX.123.456.789.vel")
}

func test(c echo.Context) error {
	//go clean -testcache
	cmd := exec.Command("go", "clean", "-testcache")

	_, err := cmd.Output()

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{err.Error()})
	}

	cmdtest := exec.Command("go", "test", "./demotest", "-v")
	// sed to pipe the output to
	sed := exec.Command("sed", "s/\\x1B\\[[0-9;]\\{1,\\}[A-Za-z]//g")

	testRes, err := cmdtest.Output()
	sed.Stdin = bytes.NewReader(testRes)

	// run the clean after the initial output
	clean, _ := sed.Output()

	if err != nil {
		return c.Blob(http.StatusBadRequest, echo.MIMETextPlain, clean)
	}

	//fmt.Println(string(testRes))
	return c.Blob(http.StatusOK, echo.MIMETextPlain, clean)
}

func getServices(c echo.Context) error {

	// c.ParamNames() to get the names and match what ones are sent
	// if actions are called
	// ParamValues() gets all the parameter values

	name := c.QueryParam("inputMRXID")
	outputID := c.QueryParam("outputMRXID")

	ins, err := transformations.ServicesMrxPathFinder(name, outputID, "")

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSONPretty(http.StatusOK, ins, "    ")
}

// autoELT converts a data type to another using the metarex register
func autoELT(c echo.Context) error {

	params, err := c.FormParams()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{"no parameters received"})
	}

	// get data to be transformed
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request().Body)
		if len(bodyBytes) == 0 {
			return c.JSON(http.StatusBadRequest, ErrorMessage{"no data received"})
		}

	} else {
		return c.JSON(http.StatusBadRequest, ErrorMessage{"no data received"})
	}

	// get the source and destination values from the register
	name := c.QueryParam("inputMRXID")
	outputID := c.QueryParam("outputMRXID")
	mapping := c.QueryParam("mapping")

	var mapEnabled bool
	if mapping == "true" {
		mapEnabled = true
	}

	// find the paths
	dataActions, err := transformations.MrxPathFinder("", name, outputID, mapEnabled)
	if err != nil {

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// set the logging for transformations
	metadata := mrxlog.NewMRX("")
	for i, dataAction := range dataActions.Actions {

		// generate the data for the run
		// without overwriting the input metadata
		outputData := make([]byte, len(bodyBytes))
		copy(outputData, bodyBytes)

		// log the action position
		metadata.LogInfo(fmt.Sprintf("Trying register path %v ", i))

		// Action history is for logging the data transformations
		var ActionHistory *mrxlog.MRXHistory
		for _, action := range dataAction {

			// make sure the logging body is set on the first run
			if ActionHistory == nil {
				ActionHistory = mrxlog.NewMRX(action.DataID(), mrxlog.WithAction(action.ActionType()), mrxlog.WithOrigin("inputFile"))
			} else {
				// else add the children
				ActionHistory = ActionHistory.PushChild(
					*mrxlog.NewMRX(action.DataID(), mrxlog.WithAction(action.ActionType())),
				)
			}
			// log what transformation is happening
			ActionHistory.LogInfo(fmt.Sprintf("Register path %v :%v", i, action.ActionType()))
			// actually transform the data

			outputData, err = action.Transform(outputData, params)

			if err != nil {
				ActionHistory.LogError(fmt.Sprintf("Register path %v :%v", i, err.Error()))
				break
				// then move to the next path if an error occurs
			}
		}

		// if success save the data and move on
		if err == nil {
			bodyBytes = outputData
			ActionHistory = ActionHistory.PushChild(*mrxlog.NewMRX(outputID))
			ActionHistory.LogInfo(fmt.Sprintf("Register path %v succesfully converted %v to %v", i, name, outputID))
			break
		} else {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	return c.Blob(http.StatusOK, dataActions.OutputFormat, bodyBytes)
}

/*

type sphericalCoords struct {
	Azimuth   float64 `yaml:"azimuth"`
	R         float64 `yaml:"distance"`
	Elevation float64 `yaml:"elevation"`
}



// metadataExtractorTransformer opens the MRX and gets the metadata.
// before attempting to transform it to the output MRXID
// this was built for the command line version of the demo
func metadataExtractorTransformer(inputFile, outputID string) error {

	//  open the input file
	f, err := os.Open(inputFile)

	noIDlog := mrxlog.NewMRX("")
	if err != nil {
		noIDlog.LogError(err.Error())
		// @TODO decide if error needs to be returned
		return err
	}

	// decode the MRX
	dataPoints, err := decode.ExtractStreamData(f)
	if err != nil {
		noIDlog.LogError(err.Error())
		return err
	}

	// for each data type in the file (not including the manifest)
	// extract it to a target format
	for _, data := range dataPoints[:len(dataPoints)-1] {

		// generate the log body
		metadata := mrxlog.NewMRX(data.MRXID, mrxlog.WithAction("Extracting"), mrxlog.WithOrigin(inputFile)) //MrxID: data.MRXID, Action: "Identify", Extra: map[string]any{"Origin": mrxFile}}
		metadata.LogInfo(fmt.Sprintf("Searching register for %v", data.MRXID))
		// find the transformations here
		dataActions, err := transformations.MrxPathFinder(inputFile, data.MRXID, outputID, true)

		if err != nil {
			metadata.LogError(err.Error())
		}

		for i, dataAction := range dataActions.Actions {

			// generate the data for the run
			outputData := make([][]byte, len(data.Data))

			for i := range outputData {
				outputData[i] = make([]byte, len(data.Data[i]))
				copy(outputData[i], data.Data[i])
			}

			// log the action position
			metadata.LogInfo(fmt.Sprintf("Trying register path %v ", i))

			// Action history is for logging the data transformations
			var ActionHistory *mrxlog.MRXHistory
			for _, action := range dataAction {
				// make sure the body is set on the first run
				if ActionHistory == nil {
					ActionHistory = mrxlog.NewMRX(action.DataID(), mrxlog.WithAction(action.ActionType()), mrxlog.WithOrigin(inputFile))
				} else {
					ActionHistory = ActionHistory.PushChild(
						*mrxlog.NewMRX(action.DataID(), mrxlog.WithAction(action.ActionType())),
					)
				}
				// log what transformation is happening
				ActionHistory.LogInfo(fmt.Sprintf("Register path %v :%v", i, action.ActionType()))
				// actually transform the data
				outputData, err = action.Transform(outputData)
				if err != nil {
					ActionHistory.LogError(fmt.Sprintf("Register path %v :%v", i, err.Error()))
					break
					// then move to the next path if an error occurs
				}
			}

			// if success save the data and move on
			if err == nil {
				data.Data = outputData
				ActionHistory = ActionHistory.PushChild(*mrxlog.NewMRX(outputID))
				ActionHistory.LogInfo(fmt.Sprintf("Register path %v succesfully converted %v to %v", i, data.MRXID, outputID))
				break
			}
		}

		// convert the strings to data type for method of printing them
		if err == nil {
			//	res := make([]json.RawMessage, len(data.Data))
			//	for i := range res {
			//		res[i] = data.Data[i]
			//	}

			f, _ := os.Create("./demodata/output.mrx")
			in := make(chan []byte, len(data.Data))

			//var inputData []any
			//b, _ := os.ReadFile("./demodata/transformedDataVelocity.mrx.json")
			//err = json.Unmarshal(b, &inputData)

			for _, dat := range data.Data {

				in <- dat
			}

			go func() {
				for {

					if len(in) == 0 {
						close(in)
						break

					}
				}
			}()

			// encode the as a single stream MRX
			demoConfig := encode.Configuration{Version: "pre alpha",
				Default:          encode.StreamProperties{StreamType: "Example XYZ data", FrameRate: "1/1", NameSpace: "https://metarex.media/reg/" + outputID},
				StreamProperties: map[int]encode.StreamProperties{0: {NameSpace: outputID}},
			}
			err := encode.EncodeSingleDataStream(f, in, demoConfig)
			fmt.Println(err)

		}

	}
	// output

	/*

		API design, for velocity and acceleration.
		Send large arrays or individual


	return nil
}
*/
