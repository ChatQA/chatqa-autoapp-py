


# Python程序编译

```
export GOOS=windows && export GOARCH=amd64
go build
```

## Windows

```shell
mkdir xxxxxxxx
cd xxxxxxxx
vim app.py
pigar generate --auto-select --question-answer=yes --index-url=http://mirrors.cloud.aliyuncs.com/pypi/simple/ && pip install -r requirements.txt && pyinstaller --onefile app.py
```

## macOS

在Linux上使用osxcross编译

## Linux