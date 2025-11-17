# ⚠️ CGO Kurulum Gerekli

## Hızlı Çözüm (5 dakika)

### Yöntem 1: WinLibs ile (Önerilen - Admin gerektirmez)

1. https://winlibs.com/ adresine gidin
2. **UCRT runtime** seçeneğini indirin (örn: `winlibs-x86_64-posix-seh-gcc-13.2.0-mingw-w64ucrt-11.0.1-r5.zip`)
3. ZIP'i `C:\mingw64` klasörüne çıkarın
4. PATH'e ekleyin:

```powershell
# Geçici olarak (bu terminal için)
$env:PATH += ";C:\mingw64\bin"

# Kalıcı olarak (System Settings)
# Windows tuşu > "environment" ara > "Environment Variables" > 
# User variables > Path > Edit > New > C:\mingw64\bin ekle
```

5. Yeni terminal açıp doğrula:
```powershell
gcc --version
```

### Yöntem 2: MSYS2 ile (Daha kapsamlı)

```powershell
# MSYS2 indirme
# https://www.msys2.org/ adresinden installer'ı indir ve çalıştır

# MSYS2 terminal'de şunu çalıştır:
pacman -S mingw-w64-ucrt-x86_64-gcc

# PATH'e ekle:
$env:PATH += ";C:\msys64\ucrt64\bin"
```

### Yöntem 3: Chocolatey (Admin gerekli)

```powershell
# PowerShell'i "Run as Administrator" ile aç
choco install mingw -y
```

## Build Sonrası

```powershell
# Doğrula
gcc --version
go env CGO_ENABLED  # 1 olmalı

# Build
cd c:\Projects\personel\go\clients
go build -o clientinfo.exe .

# Çalıştır
.\clientinfo.exe
```

## Alternatif: Web-tabanlı UI?

Eğer CGO kurmak istemiyorsanız, web tabanlı arayüz yapabiliriz (tarayıcıda açılır ama yerel çalışır). İster misiniz?
