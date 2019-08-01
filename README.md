# MySQL All databases backup with Go

### Replace database user, pass and host
```go
username := ""
password := ""
hostname := ""
dumpDir := "/opt/dumps"
backDir := "/opt/backup/"
```

```bash
[dther@opslab mysql-all-dbs-backup]$ go build .
[dther@opslab mysql-all-dbs-backup]$ mkdir /opt/backup #or change dump and backup dir in source 
[dther@opslab mysql-all-dbs-backup]$ ./mysql-all-dbs-backup
```
