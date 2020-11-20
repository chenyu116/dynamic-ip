version=$(git describe --abbrev=0 --tags)
branch=$(git rev-parse --abbrev-ref HEAD)
commit=$(git log --pretty=format:"%h" -1)
buildTime=$(date --rfc-3339=seconds | sed 's/ /T/')
verson='0.0.1'
echo "version: $version branch: $branch commit: $commit buildTime: $buildTime"
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags "-X main._version_=$version -X main._branch_=$branch -X main._commit_=$commit -X main._buildTime_=$buildTime" main.go
if [ "$?" -ne "0" ]; then
  echo "build fail"
  exit 1
fi
docker build -t ccr.ccs.tencentyun.com/astatium.com/node-dynamic-ip:v${version}-alpine-arm64 .
docker push ccr.ccs.tencentyun.com/astatium.com/node-dynamic-ip:v${version}-alpine-arm64
docker rmi ccr.ccs.tencentyun.com/astatium.com/node-dynamic-ip:v${version}-alpine-arm64