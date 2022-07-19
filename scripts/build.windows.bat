@echo off
cd ..
git rev-list -1 HEAD > tempFile && set /P gitCommit=<tempFile
git describe --tags --abbrev=0 > tempFile && set /P gitTag=<tempFile
del tempFile
go build -o api-server.exe -a -ldflags "-extldflags '-static' -X 'main.BuildTime=%DATE% %TIME:~0,-3%' -X 'main.GitCommit=%gitCommit%' -X 'main.GitTag=%gitTag%'" cmd\api\main.go
echo "api-server built."
go build -o web-server.exe -a -ldflags "-extldflags '-static' -X 'main.BuildTime=%DATE% %TIME:~0,-3%' -X 'main.GitCommit=%gitCommit%' -X 'main.GitTag=%gitTag%'" cmd\web\main.go
echo "web-server built."