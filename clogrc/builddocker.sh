## Save all the local projects
go mod vendor
# set up multi arch
docker buildx create --use --name multi-arch-builde

#docker login -u mmtristan -p "$DOCKER_PWD"
## build then run the docker
## docker buildx build -t "$TARGETamd" --push --platform linux/amd64 .
##docker build --tag mrx-elt-demo .
echo "Build & push to mmtristan/mrx-elt-demo-arm"
docker buildx build -t mmtristan/mrx-elt-demo-arm --push --platform linux/arm64 .
docker buildx imagetools inspect mmtristan/mrx-elt-demo-arm

echo "Build & push to mmtristan/mrx-elt-demo-amd"
docker buildx build -t mmtristan/mrx-elt-demo-amd --push --platform linux/amd64 .
docker buildx imagetools inspect mmtristan/mrx-elt-demo-amd
#docker buildx build -t mrx-elt-demo --push --platform linux/amd64 .
##docker build --tag mrx-elt-demo .
#docker run --publish 8080:8080 mrx-elt-demo

rm -R vendor
