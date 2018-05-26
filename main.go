package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
)

func main() {
	out, _ := os.Create("output.txt")
	defer out.Close()

	resp, _ := http.Get("http://data.gharchive.org/2015-01-01-15.json.gz")
	defer resp.Body.Close()

	uncompressed_resp, _ := gzip.NewReader(resp.Body)

	p := make([]byte, 256)
	fmt.Println(uncompressed_resp.Read(p))
	fmt.Println(uncompressed_resp.Read(p))
	fmt.Println(uncompressed_resp.Read(p))

	// io.Copy(out, resp.Body)
}
