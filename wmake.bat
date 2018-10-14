@echo off
set DIR=%~dp0
cd %DIR%

if "%~1" == "setup" goto setup
if "%~1" == "build" goto build
goto default

:setup
    net session >nul 2>&1
    if NOT %errorLevel% == 0 (
        echo Admin rights are needed to create symlinks
        exit 1
    )

    REM Initialize submodules
    echo Initializing submodules ...
    git submodule update --init --recursive

    REM Get development dependencies
    echo Getting development dependencies ...
    make get

    REM Get system library
    echo Getting system library ...
    go get -u golang.org/x/sys/windows

    REM Sync govendor dependencies
    echo Syncing govendor dependencies ...
    govendor sync

    REM Get gotools
    echo Getting gotools ...
    go get -u golang.org/x/tools/...

    goto end

:build
    cd %DIR%\cmd\wio\utils\io
    go-bindata -nomemcopy -pkg io -prefix ../../../../ ../../../../assets/...
    cd %DIR%\cmd\wio
    go build -o ..\..\bin\wio.exe -v

    goto end

:default
    echo wmake.bat [command]
    exit 1

:end

