# allbranding

```bash
# Install
url=https://github.com/taylormonacelli/allbranding/releases/latest/download/allbranding_Linux_x86_64.tar.gz
curl -fsSL $url | tar -C /usr/local/bin --no-same-owner -xz allbranding
```

usage example: get latest boilerpalte version
```log
$ allbranding query --releases-url=https://api.github.com/repos/gruntwork-io/boilerplate/releases --asset-regex='.*linux_amd64'
{"browser_download_url":"https://github.com/gruntwork-io/boilerplate/releases/download/v0.5.14/boilerplate_linux_amd64","version":"v0.5.14"}
```
