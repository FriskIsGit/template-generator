@echo off

echo Compiling for linus amd64
set GOOS=linux
set GOARCH=amd64
go build -o templater-linux

if %ERRORLEVEL% neq 0 (
    echo ERROR exit code %ERRORLEVEL%.
) else (
    echo The command succeeded.
)

@echo on