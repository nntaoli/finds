# build   
go build -ldflags="-w -s"   

# run   
## 替换内容关键字
./find-replace content --old "old xxx" --new "new xx" dirpath

## 批量重命名
./find-replace rename --old="old" --new="new char" dirpath
