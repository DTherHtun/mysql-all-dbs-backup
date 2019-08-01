package main

import (
	"archive/tar"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Database name
type Database struct {
	name string `json:"name"`
}

func main() {
	username := ""
	password := ""
	hostname := ""
	port := "3306"
	dumpDir := "/opt/dumps"
	backDir := "/opt/backup/"
	os.RemoveAll(dumpDir)
	os.Mkdir(dumpDir, 0777)
	dbs, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, hostname, port))
	if err != nil {
		panic(err.Error())
	}

	defer dbs.Close()

	showdbs, err := dbs.Query("SHOW DATABASES;")
	if err != nil {
		panic(err.Error())
	}
	for showdbs.Next() {
		var dbname Database
		err := showdbs.Scan(&dbname.name)
		if err != nil {
			panic(err.Error())
		}
		dumpFilename := fmt.Sprintf("%s/%s.sql", dumpDir, dbname.name)

		args := []string{"-h", hostname, "-u", username, "-p" + password, dbname.name}
		cmd := exec.Command("mysqldump", args...)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error stdout", err)
		}

		err = ioutil.WriteFile(dumpFilename, stdout, 0644)
		if err != nil {
			panic(err)
		}
	}
	defer showdbs.Close()
	tarit(dumpDir+"/", backDir)
	fmt.Println("========= All db backup! =======")
}

func tarit(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s-%s.tar", filename, time.Now().Format("20060102150405")))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}
