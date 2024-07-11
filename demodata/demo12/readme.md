# API Calls

localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnf&outputMRXID=generateFFmpeg&title=**Insert Paramter**

## Parameters

Key : Title

Values (Key - value):

- IET - Title:  IET_Panel
- SW - Title:  SpringWatch
- LP - Title:  LostPast
- LM - Title:  Cosmos-Laundromat

## Expected Input and Output

### Transformation type

known service

### Input

MRXID : MRX.123.456.789.rnf
Desc: Ninjs json metadata

### Outputs

ServiceID : generateFFmpeg
Desc : The ffmpeg script to make the video segements,  generated using the rnf input.
