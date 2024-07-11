// package transformations handles the traversal through the Metarex register and
// the transformations to different metadata types.
package transformations

import (
	"fmt"
	"net/url"

	"github.com/metarex-media/mrx-demo-handlers/mrxhandle/mrxlog"
	"github.com/metarex-media/mrx-demo-svc/register"
	"github.com/metarex-media/mrx-demo-svc/transformations/api"
	"github.com/metarex-media/mrx-demo-svc/transformations/mapping"
)

// Action is the method to transform metadata from one type to another.
// Transforming a single piece of metadata.
type Action interface {
	// Transform converts the a data point into a different type
	Transform(in []byte, params url.Values) ([]byte, error)
	// ActionType describes the transformation
	ActionType() string
	// DataID returns the ID of the data being transformed
	DataID() string
}

/*
MRXPathFinder searches for the path between two Metarex IDs, or a MetarexID and service.
returning an array of actions to transform between them. In the instance the first
series of transformations fails try the second one and so on, if there is more than one set
of transformations that is.

The current transformation paths are:

  - API transformations - where an API returns the metadata, this can take several steps to get to the destination format.
  - Mapping transformations - a one step best guess translation into the destination metadata format.
*/
func MrxPathFinder(inputOrigin, sourceID, destinationID string, mapTransforms bool) (*dataProcess, error) {

	inReg, ok := register.GetRegEntry(sourceID)
	if !ok {
		return nil, fmt.Errorf("no register entry found for %v", sourceID)
	}
	// check the output register at the last moment
	// because it may be a serivce with no set data type
	outReg, ok := register.GetRegEntry(destinationID)
	logBody := *mrxlog.NewMRX(sourceID, mrxlog.WithOrigin(inputOrigin))

	if mapTransforms {
		if !ok {
			return nil, fmt.Errorf("no register entry found for %v", destinationID)
		}
		// do a mapping transformation and hope for the best
		mappingAction := mapping.MappingAction{OutputSchema: outReg.Mrx.Spec,
			MrxID:       destinationID,
			InputFormat: inReg.MediaType, OutputFormat: outReg.MediaType,
			InputTiming: inReg.Timing,
			Mapping:     *outReg.Mrx.Mapping,
		}
		logBody.LogInfo("No direct path found, trying mapping transformation")
		return &dataProcess{Actions: [][]Action{{&mappingAction}}, OutputFormat: outReg.MediaType}, nil
		// utilise thesaurus action
	}

	// search the IDs
	// basic search of just the APIs to start
	searcher := search{}

	searcher.RegisterDive(logBody, destinationID, sourceID, []path{})

	// for each valid path
	chosenPaths := searcher.validPaths
	actions := make([][]Action, len(chosenPaths))

	// log chosen path
	for i, chosenPath := range chosenPaths {
		actions[i] = make([]Action, len(chosenPath))
		for j, vp := range chosenPath {

			target, _ := register.GetRegEntry(vp.Id)
			service := target.Mrx.Services

			// translate the parameters into API required parameters
			outParams := make([]api.Parameter, len(service[vp.array].Parameters))
			for i, param := range service[vp.array].Parameters {
				outParams[i] = api.Parameter(param)
			}
			//	fmt.Println(service, len(actions), len(service))
			actions[i][j] = &api.ApiAction{API: service[vp.array].API,
				MrxID: vp.Id, ResponseMIMEType: target.MediaType, APISchemaLocation: service[vp.array].Spec,
				APIParams: outParams,
			}
		}
	}
	if len(actions) == 0 && !mapTransforms {
		return nil, fmt.Errorf("no path found, please try again by enabling mapping with \"mapping=true\"")
	}

	// to get media types and not output types
	output := searcher.outputType
	if reg, ok := register.GetRegEntry(searcher.outputType); ok {
		output = reg.MediaType
	}
	// return some actions
	return &dataProcess{Actions: actions, OutputFormat: output}, nil
	//return actions, nil
}

type serviceInformation struct {
	APICall         string
	Description     string `json:"description,omitempty"`
	ServiceID       string
	MRXRegisterPath []string             `json:"MRXRegisterPath"`
	Params          []register.Parameter `json:"parameters"`
}

/*
ServicesMrxPathFinder checks for services in an MRX based on output types.
It only checks the services of the sourceID and is not yet recursive
e.g. metadata to image/png

It returns the serivceID, description and url of the API to make the post request to.
*/
func ServicesMrxPathFinder(sourceID, serviceType, inputOrigin string) ([]serviceInformation, error) {

	//inReg, ok := register[sourceID]
	//	if !ok {
	//		return nil, fmt.Errorf("no register entry found for %v", sourceID)
	//	}

	found := make([]serviceInformation, 0)

	searcher := search{}

	logBody := *mrxlog.NewMRX(sourceID, mrxlog.WithOrigin(inputOrigin))
	searcher.RegisterDive(logBody, serviceType, sourceID, []path{})

	for _, paths := range searcher.validPaths {
		var servi register.Services
		regPaths := []string{}
		params := []register.Parameter{}
		for _, p := range paths {
			serviID, _ := register.GetRegEntry(p.Id)
			servi = serviID.Mrx.Services[p.array]

			regPaths = append(regPaths, p.Id)
			params = append(params, servi.Parameters...)
		}

		service := serviceInformation{APICall: fmt.Sprintf("localhost:8080/autoelt?inputMRXID=%v&outputMRXID=%v", sourceID, servi.ServiceID),
			Description: servi.Description, ServiceID: servi.ServiceID, MRXRegisterPath: regPaths}

		if len(params) != 0 {
			service.Params = params
		}

		found = append(found, service)
	}
	/*
		for _, servi := range inReg.Mrx.Services {
			if servi.Output == serviceType {
				found = append(found, serviceInformation{APICall: fmt.Sprintf("localhost:80/autoelt?inputMRXID=%v&outputMRXID=%v", sourceID, servi.ServiceID),
					Description: servi.Description, ServiceID: servi.ServiceID, Paths: []string{sourceID}})
			}
		}*/

	outReg, ok := register.GetRegEntry(serviceType)
	if ok {
		// search for a route through the paths
		// and suggest mapping if the dest has a mapping map
		if outReg.Mrx.Mapping != nil {
			if outReg.Mrx.Mapping.MappingDefinitions != nil {
				found = append(found, serviceInformation{APICall: fmt.Sprintf("localhost:8080/autoelt?inputMRXID=%v&outputMRXID=%v&mapping=true", sourceID, serviceType),
					Description: fmt.Sprintf("Generically mapping %v to %v", sourceID, serviceType), ServiceID: serviceType, MRXRegisterPath: []string{sourceID}})
			}
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("no services with the type %v were found", serviceType)
	}

	// return a description
	// completeMRXPath

	return found, nil
}

type dataProcess struct {
	Actions      [][]Action
	OutputFormat string
}

// RegisterDive recursively searches the register for a path between the input and output metadata ID
// the current path is started with an empty parent array e.g. []path{}
func (s *search) RegisterDive(mrxPath mrxlog.MRXHistory, endId, currentID string, currentPath []path) {

	if Contains(currentPath, currentID) {
		return // return to prevent recursive register searhes
	}

	// MRX log each search in debug saying yes no, where to next
	APICallsMRX, _ := register.GetRegEntry(currentID)
	APICalls := APICallsMRX.Mrx.Services

	// if dead end

	if len(APICalls) == 0 {
		mrxPath.LogDebug(fmt.Sprintf("No path found following %v", currentPath))
	}

	for i, call := range APICalls {
		mrxChild := mrxPath.PushChild(*mrxlog.NewMRX(call.ID, mrxlog.WithAction("Searching the register")))
		// either find it or keep on searching
		newPath := append(currentPath, path{Id: currentID, array: i})

		if call.ID == endId || call.ServiceID == endId || call.Output == endId {

			// assign the path
			mrxChild.LogInfo(fmt.Sprintf("Found register path %v", newPath))
			s.validPaths = append(s.validPaths, newPath)
			if s.outputType == "" {
				s.outputType = call.Output
			}
			// return leave the return
		} else {

			mrxChild.LogDebug(fmt.Sprintf("building register search path: %v", newPath))
			s.RegisterDive(*mrxChild, endId, call.Output, newPath)
		}
	}
}

// Contains searches a path to see if it contains an ID,
// to prevent endless loops
func Contains(paths []path, ID string) bool {
	for _, path := range paths {
		if path.Id == ID {
			return true
		}
	}

	return false
}

type search struct {
	outputType string
	depth      int
	validPaths [][]path
}

type path struct {
	Id    string
	array int
}
