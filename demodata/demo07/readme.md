# API Calls

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fill
localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fault
localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=stagger

## Expected Input and Output

### Transformation type

known service

### Input

MRXID : MRX.123.456.789.bat
Desc : The Riedel battery format

### Outputs

ServiceID : fill
Desc : A scaled PNG showing the battery percentage

ServiceID : fault
Desc : A JPEG showing a visualisation of the battery Error code.

ServiceID : fill
Desc : A staggered PNG giving a rough idea of the battery percentage
