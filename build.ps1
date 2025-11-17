# Windows GUI uygulaması olarak derle (konsol penceresi olmadan)
go build -ldflags "-H windowsgui" -o clients.exe .\internal\.

Write-Host "Derleme tamamlandi: clients.exe" -ForegroundColor Green
