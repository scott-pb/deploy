@CHCP 65001
@echo off

:: Build args
if %time:~0,2% leq 9 (set hour=0%time:~1,1%) else (set hour=%time:~0,2%)
SET datetime=%date:~3,4%%date:~8,2%%date:~11,2%-%hour%%time:~3,2%
SET CGO_ENABLED=0
SET GOOS=linux
SET ROOT=%cd%
mkdir %ROOT%\bin\
@echo Start build(v%datetime%) running 🚀🚀 ...

go build -o .\bin\deploy -gcflags=all="-N -l" -ldflags="-X main.version=v%datetime%" -trimpath
@echo soga_admin ✔️
pause
exit