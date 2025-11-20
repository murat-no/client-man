windows versiyon bilgileri için <br>
goversioninfo -64 -o resource.syso versioninfo.json <br>

windows gui derlemek için <br>
go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\. <br>


//goversioninfo -64 -o resource.syso versioninfo.json
//go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\.
//go build -trimpath -ldflags="-s -w -H windowsgui" -o client-manager.exe .\internal\.
//go run .\internal\. 2>&1