@REM @echo off

set BIN_DIR=C:\tools\bin
set FILES_DIR=C:\ProgramData\ArquivosRecebidos
set PORT=10777

mkdir "%BIN_DIR%"
mkdir "%FILES_DIR%"
copy receivefiles.exe %BIN_DIR%\receivefiles.exe

sc create ReceiveFiles ^
    binPath="\"C:\tools\bin\receivefiles.exe\" service --port %PORT% --save-to=\"%FILES_DIR%\"" ^
    start=delayed-auto
sc start ReceiveFiles