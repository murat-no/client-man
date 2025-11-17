@echo off
echo CGO etkin sekilde derleniyor...
set CGO_ENABLED=1
go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\.
echo Derleme tamamlandi: client-manager.exe
