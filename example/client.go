package main

import (
	"io"
	"log"
	"mylib/sftp"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	client, err := sftp.NewSFTPClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	server := os.Args[1]
	if err := client.Connect(server); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	name := os.Args[2]
	client.GetFile(name);
	res, err := sftp.ReadResponse(client.Conn)
	fmt.Println(res)
	if err != nil {
			return 
	}
	wFile, err := os.Create("download/" + filepath.Base(name))
	if err != nil {
			log.Fatal(err) 
	}
	defer wFile.Close()

	if res.Code() == sftp.ResponseGetOk {
		io.Copy(wFile, client.Conn)
	}

}
