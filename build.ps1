# Windows GUI uygulaması olarak derle (konsol penceresi olmadan)
go build -ldflags "-H windowsgui" -o client-manager.exe .\internal\.

Write-Host "Derleme tamamlandi: client-manager.exe" -ForegroundColor Green
