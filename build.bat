@echo off
echo CGO etkin sekilde derleniyor...
set CGO_ENABLED=1
go build -ldflags "-H windowsgui" -o clients.exe .\internal\.
echo Derleme tamamlandi: clients.exe
