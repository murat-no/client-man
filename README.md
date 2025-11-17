windows versiyon bilgileri için <br>
goversioninfo -64 -o resource.syso versioninfo.json <br>

windows gui derlemek için <br>
go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\. <br>
