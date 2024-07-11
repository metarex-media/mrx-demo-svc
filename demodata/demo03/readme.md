# API Calls

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.gpx&outputMRXID=toW3C

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.wav&outputMRXID=ToWaveform

## Expected Input and Output

### Transformation type

Known service

### Input

MRXID : MRX.123.456.789.gpx
Desc : The gps exchange format

MRXID : MRX.123.456.789.wav
Desc : A wav audio file of some bird song

### Outputs

MRXID : MRX.123.456.789.gps
Desc : The w3c gps browser format

ServiceID : ToWaveform
Desc : A png visualisation of the wav file

Gpx samples from : [github](https://github.com/gps-touring/sample-gpx/tree/master)
