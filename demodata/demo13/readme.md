# API Calls

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnf&outputMRXID=MXFToGraph
localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnf&outputMRXID=MXFToReport

## Parameters

## Expected Input and Output

### Transformation type

known service

### Input

MRXID : MRX.123.456.789.mxf
Desc: AN MXF file, with hopefully some ISXD data

### Outputs

ServiceID : MXFToGraph
Desc : Convert an mxf file to a graph detailing its test report

ServiceID : MXFToReport
Desc : Run a series of tests on the mxf file, with the report given as a yaml
