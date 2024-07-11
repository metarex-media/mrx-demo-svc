# mrx-demo-svc

| [ETL](#etl) | [build](#bld) | [transform](#xfm) | [demos](#demos)
| [1](#d1) | [2](#d2) | [3](#d3) | [4](#d4) | [5](#d5) | [6](#d6)
| [7](#d7) | [8](#d8) | [9](#d9) | [10](#d10) | [11](#d11) | [12](#d12) |

Welcome to `mrx-demo-svc` a web service for performing **ETL** (**E**xtract
**T**ransfom **L**oad) operations in the MetaRex Website demos. To find out
more about MetaRex and our mission, visit the [MetaRexwebsite][01].

Follow the 12 [demos](#demo-breakdown) to learn how MetaRex can make it easy to
integrate new and unknown data into your workflows.

All these demos utilize the [MetaRex register][02], to drive the single generic
ETL pipeline that every demo uses. Instead of having to make a pipeline for
every input and output transformation, you can now use the same pipeline for
every use case. This means less time coding about the metadata and more using
the metadata.

<a id= "etl"></a>

## What is ETL (and what is ELT)

**ETL** is a data extraction method and it stands for **E**xtract **T**ransform
**L**oad and is the key to the way Metarex works. The services in this repo do
the following:

* **E**xtract - extract the source data from a container or via an API and
  fetch the corresponding MetaRex register entry
* **T**ransform - into a known target data type using fixed algorithms or
   published services or via schema inspection logic using information found in
   the MetaRex Register entry.
* **L**oad - load the data in its target format into a workflow

There is an interesting description of the ETL vs ELT difference on the [AWS
website][04] where they come to the conclusion that ELT is the norm when you
assume a large volume of data is to be converted between databases with custom
code; often where the format of the data is quasi-static over long periods.

MetaRex assumes that there is a large volume of different and changeable data
types because the media industry is constantly innovating and thus well
described metadata types improve interoperability and allow easier data
versioning and more business flexibility.

<a id="bld"></a>

## Building & running the demos

These services are intended to be deployed in a docker container and the
endpoints of the services are listed in a local version of the MetaRex
register. The demos are active on the [metarex.media][01] website.

### Pre requisites

#### libxml2

[go-xsd-validate](https://godoc.org/github.com/terminalstatic/go-xsd-validate)
requires libxml2 to run. This step is only required if you want to run the
demo straight from the command line and not in the docker container.

Install `libxml2` via distribution package manager or from source, e.g.

```sh
sudo apt-get install libxml2 xsltproc
# use xsltproc to display the xmllib2 version (2.09.14 at time or writing)
xsltproc --version
```

#### golang

Install golang (see [go.mod](go.mod) for the required version) from the
[official source](https://go.dev/doc/install)

#### Docker

Install [Docker](https://docs.docker.com/engine/install/)

### Building and running the docker image

Run the following commands to build the self contained demo image using the Dockerfile.
This exposes a service on port 8080 that does all the auto ETL transformations.

```cmd
docker build --tag mrx-demo-svc .docker
docker run --publish 8080:8080 mrx-demo-svc
```

Please note if you are cross compiling your docker image for arm and amd, then
you can do this with [buildx][bx]. This boils down to the following commands.

```cmd
# set cross compiler
docker buildx create --use --name multi-arch-builder
# build the locally cross compiled image
docker buildx build --tag mrx-demo-svc-arm --load --platform linux/arm64 .

# or  build and push to docker hub
# login to docker using a Personal Access Token (PAT)
docker login -u $DOCKER_USR -p "$DOCKER_PAT"
# push to the remote repo
docker buildx build -t $DOCKER_USR/mrx-demo-svc-arm --push --platform linux/arm64 .
```

These docker images should be available on the [MetaRex Docker Hub][md].

### Running from the command line

If you want to build and run the demo as a command line tool then run the
following steps.

This will build a services api on localhost:9000 and the auto ETL server on
localhost:8080.

```cmd
# build the auto ETL server
go build -o mrx-demos
# build the services API
cd api && go build api || cd ..

# Run both servers using the run script
chmod u+x ./clogrc/run.sh && ./clogrc/run.sh
```

Please note if you are cross compiling for other operating systems, then you can
[cross compile with native go][cc].

### Testing the server is running

After building the demo, ensure the image is running by calling the server test.

`curl localhost:8080/test`

or by visiting [localhost:8080/test](http://localhost:8080) in the browser

This will give a go formatted output of pass or fail. Pass means the
transformations are all running as intended.

You can read the inline docs by running the godocs command:

```sh
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=localhost:6060 -goroot=.
```

### Methodology

All data types have a **`metarexId`**  (`mrxId`) that identifies the class of
that data. To perform a transform, the code looks at the `mrxId` for the
transform source and destination. The MetaRex Register is then searched to
entries with the source `mrxId` containing explicit services to convert to the
`mxrId` or `Content-Type` for the destination. This results in a list of transformations to
convert the data into the desired format.

The current transformation paths between register entries are:

* **API transformations** - where an explicit API service transforms the
metadata into another `mrxId` type or `Content-Type`. Several API
transformations may be called to convert to the destination format.
* **Heuristic Mapping transformations** - a one step translation to the
destination format based on mapping the terms in the source schema to the terms
in the destination schema.

### Logging

mrxlog is used to track the recursive nature of the register, offering the full
path taken through the MRX register when making transformations. As well as the
actions that led to that log being made, such as `mapping` or `calling service
at example.com/service`.

A log entry will look similar to

```cmd
2024-04-23T10:33:41Z INF mrxlog/mrxlog.go:242 Register path 0 succesfully converted MRX.123.456.789.rnc to MRX.123.456.789.rnf MRXPath="{MRX.123.456.789.rnf  {MRX.123.456.789.rnf mapping to MRX.123.456.789.rnf <nil> map[] inputFile} map[] }" chainID=6c450a8f8f8d8835 parentID=961ec17a01baa198
```

The following properties are used for tracking transformation paths:

* ChainID is the xxh64 hash of all the MRXIDs that have led to this
transformation. e.g. A chain of `MRX.123.456.789.rnc` then
`MRX.123.456.789.rnf` will have a chainID of `6c450a8f8f8d8835`.
* ParentID is the ChainID of the parent MRXID, for the  chainID of
`6c450a8f8f8d8835` there will be a ParentID of `961ec17a01baa198`, as this is
just the hash of `MRX.123.456.789.rnc`

It can be saved as a database or straight to the command line. Error ids can
then be searched to find elements in the chain that caused a problem. The
example SQL below assumes the logs have been saved to sqlite or similar SQL
database for processing.

```sql
SELECT * FROM log WHERE chainId="6c450a8f8f8d8835"
```

<a id="xfm"></a>

## Transformations

### API transforms

API transformations POST the data to a service, where it is transformed into
another data type. This service is found by inspecting the MetaRex Register
for entries with an `mrxId` that matches the source data.

### Mapping transforms

Mapping is a best guess data translation, using the schema for the target as a
blueprint to find a mapping from the source schema.

It uses the following simple algorithm (that you could extend or add some AI to
potentially improve its accuracy):

* flattens the input data.
* it finds the flat destination layout from the schema.
* it matches the source field names with any destination field names.
* if multiple destination field names are found then the matching is based on
  similar depth of the field in the schema's hierarchy.
* Fields are handled in alphabetical order, to preserve the array order of
  objects.

Default behaviours:

* 1 to 1 mapping of source to destination fields.
* objects are built as they go, no object is transformed from one to another
  due to the flat nature of the mapping.
* fields which aren't translated are discarded, can be optionally saved as an
  object with a `MissedTags` key.
* mapping dictionary is expected field key with an array of potential field
  names.

Default transformations type:

* non-numbers are transformed as strings e.g. `fmt.Sprintf("%v", data)`
* both integers and floats become floats (i.e. JSON numbers).
* floats are floor() to ints, ints are not floor().
* boolean values are not transformed
* single values are appended to an array, of the same type.
* arrays types follow the same rules as the single values. e.g. floats are
  floor() into an integer array.

Please Note this is a demo repo to show how MetaRex works. Not every mapping
will is ready for production, especially for more complicated data types like
2d arrays. If you would like to sponsor **algorithm improvements** to make the
heuristic mapper ready for your specific production needs, we are looking for
[sponsors][01c] to fund this work. Please [get in touch][01c].

<a id="demos"></a>

## Demo Breakdown

Each demo uses a single service endpoint to perform a transformation:

```url
localhost:8080/autoETL?inputMRXID=<mrxSId>&outputMRXID=<mrxTId | svcId>[&mapping=true]
```

<a id="mapping"></a>

* `<mrxSId>` is the metarexId of the **s**ource data.
* `<mrxTId | svcId>` is either:
  * `mrxTId` the metarexId of the **t**arget data type to invoke the heuristic
  mapping algorithm
  * `svcId` the serviceId in the metarex register entry for the source data.
* `mapping=true` and optional parameter that forces heuristic mapping

For the sake of the demo a mapping parameter has been added to demonstrate the
transformation.

Any additional parameters are passed to the handler for a specific `svcId`. e.g.
in [demo 12](#d12) a parameter called title is required.

<a id="d1"></a>

### Demo 1 - Signiant Jet

This demo takes [ninjs metadata][d1n] and transforms it into one of the
following formats:

* [a simple markdown json][d1j] This is done via a service transformation and
this format is then used to generate mark down files.

e.g. the ninjs file [ap_audio.json][d1a] is converted with a POST request to a
simple markdown file that can be used in a GUI:

```url
localhost:8080/autoETL?inputMRXID=MRX.123.456.789.njs&outputMRXID=toNewsMD
```

returns

```json
{
    "title": "# Italy: Italy Demonstration of the workers congresses and conferences",
    "shortSummary": "Demonstration of the workers of congresses and conferences against he measures taken by the Government that closed all the sporting centers to contrast the Covid-19 emergency.\nRome (Italy), October 27th 2020\nPhoto Samantha Zucchi /Insidefoto/Sipa USA)",
    "slug": ""
}
```

<a id="d2"></a>

### Demo 2

[How does MetaRex work](https://metarex.media/docs/)

This interactive demo is still work in progress.

<a id="d3"></a>

### Demo 3

This demo takes [gpx coordinate metadata][d3g] and transforms it into one of the
following formats:

* [w3c coordinates][d3x] This is done via a service transformation.

e.g. the data file [Newhaven_Brighton.gpx][d3d] converted with a POST request to

```url
localhost:8080/autoETL?inputMRXID=MRX.123.456.789.gpx&outputMRXID=toW3C
```

and returns this json

```json
[
    {
        "coords": {
            "latitude": 50.784192,
            "longitude": 0.058557,
            "altitude": 6
        }
    }
    ...
    {
        "coords": {
            "latitude": 50.82299,
            "longitude": -0.15617,
            "altitude": 13.838
        }
    }
]
```

This demo also takes [wav audio files][d3w] and transforms them into one of the
following formats:

* A PNG visualisation of the audio via a service transformation.

e.g. the wav file [European Robin - short.wav][d3f] is converted with a POST
request to

```url
localhost:8080/autoETL?inputMRXID=MRX.123.456.789.wav&outputMRXID=ToWaveform
```

and returns the png:

![A sound wave of a Robin][d3p]

<a id="d4"></a>

### Demo 4

This demo takes [yaml tagged metadata][d4y],[json tagged metadata][d4j] or [csv
tagged metadata][d4c] and transforms it into the following formats:

* [Responsive narrative factory csv](./register/registerEntries/MRX.123.456.789.rnf.json)
This is done by generically [mapping](#mapping) the input into the rnf format.

e.g. [IET.yaml](./demodata/demo04/IET.yaml)
is converted with a POST request to `localhost:8080/autoETL?inputMRXID=MRX.123.456.789.rny&outputMRXID=MRX.123.456.789.rnf&mapping=true`.
when using json the request becomes `localhost:8080/autoETL?inputMRXID=MRX.123.456.789.rnj&outputMRXID=MRX.123.456.789.rnf&mapping=true`
and when using csv it become `localhost:8080/autoETL?inputMRXID=MRX.123.456.789.rnc&outputMRXID=MRX.123.456.789.rnf&mapping=true`.
The addition parameter of `mapping=true` is required for every request.

```csv
{
chapter,in,metadataTags,out,storyline-importance
logo,0,"{""Speaker"":[""Rosana Prada""],""segment"":"""",""topics.primary"":""intro"",""topics.secondary"":"""",""topics.shot"":""intro"",""topics.subject"":""sting"",""warnings.predation"":false,""warnings.procreation"":false,""warnings.threat"":false}",167,10
chapter001,168,"{""Speaker"":[""Rosana Prada""],""segment"":"""",""topics.primary"":""intro"",""topics.secondary"":"""",""topics.shot"":""Rosanna"",""topics.subject"":""intro"",""warnings.predation"":false,""warnings.procreation"":false,""warnings.threat"":false}",492,3
...
chapter002,70808,"{""segment"":"""",""topics.primary"":""Next Gen Audio"",""topics.secondary"":""loudness"",""topics.shot"":""logo"",""topics.subject"":""close"",""warnings.predation"":false,""warnings.procreation"":false,""warnings.threat"":false}",70808,10

}
```

<a id="d5"></a>

### Demo 5

This demo takes [Venera QC report](./register/registerEntries/MRX.123.456.789.vqc.json) and transforms it into one of the following
formats:

* A PNG visualisation of the errors and warnings.
This is done via a service transformation.

e.g. [PulsarReport2Fail.xml](./demodata/demo05/PulsarReport2Fail.xml)
is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.vqc&outputMRXID=tograph`

![A graph representation of PulsarReport2Fail.xml](./demotest/testdata/expected/PulsarReport2Fail.png "Pulsar Report2 Fail")

<a id="d6"></a>

### Demo 6

This demo takes [unreal engine camera metadata](./register/registerEntries/MRX.123.456.789.ghj.json)
and transforms it into the following formats:

* [Demo XML camera format](./register/registerEntries/MRX.123.456.789.cxm.json)
This is done by generically [mapping](#mapping) the input into the xml format.

e.g. [beach_camera.json](./demodata/demo06/beach_camera.json)
is converted with a POST request to `localhost:8080/autoETL?inputMRXID=MRX.123.456.789.ghj&outputMRXID=MRX.123.456.789.cxm&mapping=true`.
The addition parameter of `mapping=true` is required for every request.

```xml
<camera xmlns="">
    <cameraItem>
        <intrinsicMatrix>5440.000226592928</intrinsicMatrix>
        <intrinsicMatrix>0</intrinsicMatrix>
    ...
        <translation>0</translation>
    </cameraItem>
</camera>
```

<a id="d7"></a>

### Demo 7

This demo takes [battery metadata](./register/registerEntries/MRX.123.456.789.bat.json) and transforms it into one of the following
formats:

* A PNG visualisation of the battery percentage.
This is done via a service transformation.

e.g. [good.json](./demodata/demo07/good.json)
is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.bat&outputMRXID=fill`

![A battery percentage visualisation filled with green](./demotest/testdata/expected/good_out.png "Battery Status")

* A PNG visualisation of any errors.
This is done via a service transformation.

e.g. [flame.json](./demodata/demo07/fire.json)
is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.bat&outputMRXID=fault`

![A cartoon flame](./demotest/testdata/expected/fire_out.jpg "Fire!")

<a id="d8"></a>

### Demo 8

This demo takes [jpeg images with C2PA metadata](./register/registerEntries/MRX.123.456.789.jpg.json) and transforms it into one of the following
formats:

* The C2PA header metadata present in the file.
This is done via a service transformation.

e.g. <img src="./demodata/demo08/truepic-20230212-library.jpg" alt="A Library" width="300"/>

is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.jpg&outputMRXID=extractHeaderc2pa`

```json
{
  "active_manifest": "com.truepic:urn:uuid:3e3ecbb7-0fa8-44dc-a2ad-f211bbc0012b",
  "manifests": {
    "com.truepic:urn:uuid:3e3ecbb7-0fa8-44dc-a2ad-f211bbc0012b": {
      "claim": {
    ...
        }
    }
}
```

<a id="d9"></a>

### Demo 9

[Here?](https://metarex.media/app-nab2024/demo09)

<a id="d10"></a>

### Demo 10

This demo takes [ninjs metadata](./register/registerEntries/MRX.123.456.789.njs.json) and transforms it into one of the following
formats:

* [NewsML](./register/registerEntries/MRX.123.456.789.nmj.json)
This is done via a service transformation.

e.g. [ntb_text.json](./demodata/demo10/ntb_text.json)
is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.njs&outputMRXID=toNewsML`

```xml
<newsItem xmlns="http://iptc.org/std/nar/2006-10-01/" guid="" version="11" standard="NewsML-G2" standardversion="2.24">
    <rightsInfo>
        <copyrightHolder uri="">
            <name>NTB</name>
        </copyrightHolder>
        <copyrightNotice></copyrightNotice>
    </rightsInfo>
    ...
</newsItem>
```

<a id="d11"></a>

### Demo 11

Get Rexy's HandID metadata and convert it to another movie database format.

<a id="d12"></a>

### Demo 12

This demo takes [rnf csv](./register/registerEntries/MRX.123.456.789.rnf.json) and transforms it into one of the following
formats:

* An ffmpeg script to make the rnf media segments.
This is done via a service transformation.

e.g. [IET.csv](./demodata/demo12/IET.csv)
is converted with a POST request to `http://localhost:8080/autoETL?inputMRXID=MRX.123.456.789.rnf&outputMRXID=generateFFmpeg&title=IET`.
Please ensure the title parameter is used to identify the media type.
The available titles are:

* LP - lostpast
* SW - springwatch
* LM - Cosmos-Laundromat
* IET - IET

```text
ffmpeg -y -i ./rnf/bbb/audio.wav -ss 00:00:00.000 -t 00:00:06.680 -acodec pcm_s16le -ac 1 -ar 16000 ./rnf/segment0000.wav
ffmpeg -y -start_number 0 -framerate 25 -i ./rnf/bbb/frame%05d.jpg -i /segment0000.wav -frames:v 167 -vcodec mpeg4 -r 25 -q:v 0 ./rnf/IET_Panel_0000_Rosana-Prada,intro,intro,sting.mp4
...
ffmpeg -y -i ./rnf/bbb/audio.wav -ss 00:47:12.280 -t 00:47:12.320 -acodec pcm_s16le -ac 1 -ar 16000 ./rnf/segment0082.wav
ffmpeg -y -start_number 70808 -framerate 25 -i ./rnf/bbb/frame%05d.jpg -i /segment0082.wav -frames:v 70808 -vcodec mpeg4 -r 25 -q:v 0 ./rnf/IET_Panel_0082_Next-Gen-Audio,loudness,logo,close.mp4

```

[01]:  https://metarex.media
[01c]: https://metarex.media/contact/
[02]:  https://www.metarex.media/ui/reg/
[04]:  https://aws.amazon.com/compare/the-difference-between-etl-and-elt/
[bx]:  https://www.docker.com/blog/how-to-rapidly-build-multi-architecture-images-with-buildx/
[md]:  https://hub.docker.com/metarexmedia/
[cc]:  https://rakyll.org/cross-compilation/

[d1a]: ./demodata/demo01/ap_audio.json
[d1n]: ./register/registerEntries/MRX.123.456.789.njs.json
[d1j]: ./register/registerEntries/MRX.123.456.789.nmd.json

[d3g]: ./register/registerEntries/MRX.123.456.789.gpx.json
[d3x]: ./register/registerEntries/MRX.123.456.789.gps.json
[d3d]: ./demodata/demo03/Newhaven_Brighton.gpx
[d3w]: ./register/registerEntries/MRX.123.456.789.wav.json
[d3f]: ./demodata/demo03/European%20Robin%20-%20short.wav
[d3p]: <./demotest/testdata/expected/European Robin - short_out.png> "Robin Soundwave"

[d4y]: ./register/registerEntries/MRX.123.456.789.rny.json
[d4j]: ./register/registerEntries/MRX.123.456.789.rnj.json
[d4c]: ./register/registerEntries/MRX.123.456.789.rnx.json
