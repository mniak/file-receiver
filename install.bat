mkdir C:\tools
mkdir C:\tools\bin
copy receivefiles.exe C:\tools\bin\receivefiles.exe

echo sc create ReceiveFiles ^
    binPath="\"C:\tools\bin\receivefiles.exe" service --port 10888" ^
    start=delayed-auto
echo sc start ReceiveFiles