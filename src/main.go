package main

import (
	"archive/zip"
	"embed"
	"flag"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FileInfo struct {
	Name string
	Size string
	Path string
}
//go:embed templates/index.tmpl
var content embed.FS
var BASE_DIR  string

func main() {
	cwd, _ := os.Getwd()

	flag.Usage = Usage
	host := flag.String("host", "localhost", "the server host")
	port := flag.Int("port", 8080, "the server port")
	directory := flag.String("directory", "", "the root directory to serve files from")
	expose := flag.Bool("expose", false, "expose the server to the local network")
	help := flag.Bool("help", false, "display help message")

	flag.StringVar(host, "l", "localhost", "the server host (short)")
	flag.IntVar(port, "p", 8080, "the server port (short)")
	flag.StringVar(directory, "d", "", "the root directory to serve files from (short)")
	flag.BoolVar(expose, "e", false, "expose the server to the local network (short)")
	flag.BoolVar(help, "h", false, "display help message (short)")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	BASE_DIR = cwd+"/"+*directory
	r := gin.Default()

	fs := http.FileServer(http.Dir(BASE_DIR))

	r.GET("/", func(c *gin.Context){
		listHandler(fs, c, ".")
	})

	r.GET("/download/*dlDir", downloadFileOrDir)

	r.GET("/:dir/*subdir", func(c *gin.Context){
		dir := c.Param("dir")
		subdir := c.Param("subdir")
		if !isDirectory(BASE_DIR+"/"+dir){
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		listHandler(fs, c, dir+subdir)

	})

	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	log.Printf("Starting file server on %s...\n", addr)

	if *expose {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					log.Printf("Serving on http://%s:%d (local network)\n", ipnet.IP.String(), port)
				}
			}
		}
	}

	err := r.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadFileOrDir(c *gin.Context){
	dlDir:=c.Param("dlDir")

	dir, _ := filepath.Abs(BASE_DIR+"/"+dlDir)
	fileName, _ := filepath.Rel(BASE_DIR, dir)
	if !isDirectory(dir){
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName )
		c.Header("Content-Type", "application/octet-stream")
			c.File(dir)
			return

		}
		c.Writer.Header().Set("Content-Type", "application/zip")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName+".zip")
		zipWriter := zip.NewWriter(c.Writer)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name, err = filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}


			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}

			return nil
		})

		err = zipWriter.Close()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

}

func listHandler(fs http.Handler, c *gin.Context, dir string) {
	var files []FileInfo
	abs_dir, err := filepath.Abs(BASE_DIR+"/"+dir)
	if err!=nil{
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	valid, err:=exists(abs_dir); 
	if !valid{
		log.Println("not exists")
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	abs_stat, _ := os.Stat(abs_dir)
	files = append(files, FileInfo{
		Name: "..",
		Size: "-",
		Path: "..",
	})
	filepath.Walk(abs_dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, FileInfo{
				Name: info.Name(),
				Size: humanizeSize(info.Size()),
				Path: path,
			})
		}
		if info.IsDir() {
			rel_path, err := filepath.Rel(BASE_DIR, path)
			if err!=nil{
				return err
			}
			if !os.SameFile(info, abs_stat) {
				files = append(files, FileInfo{
					Name: info.Name(),
					Size: "-",
					Path: rel_path,
				})

				return filepath.SkipDir
			}
		}

		return nil
	})

	t, err := template.ParseFS(content, "templates/index.tmpl")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = t.Execute(c.Writer, files)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}


func init() {
	binding.EnableDecoderUseNumber = true
	gin.SetMode(gin.ReleaseMode)
}


