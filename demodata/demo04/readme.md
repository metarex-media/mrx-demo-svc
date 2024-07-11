# API Calls

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnj&outputMRXID=MRX.123.456.789.rnf&mapping=true
localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnc&outputMRXID=MRX.123.456.789.rnf&mapping=true
localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rny&outputMRXID=MRX.123.456.789.rnf&mapping=true

## Expected Input and Output

### Transformation type

Known service

### Input

MRXID : MRX.123.456.789.rnc
Desc : A csv tagging of media

MRXID : MRX.123.456.789.rnj
Desc : A json tagging of media

MRXID : MRX.123.456.789.rny
Desc : A yaml tagging of media

### Outputs

MRXID : MRX.123.456.789.rnf
Desc : The rnf csv format
