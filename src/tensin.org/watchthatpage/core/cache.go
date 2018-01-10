package core

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

// Minify will removes useless entries in the HTML content
func Minify(content []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	content, err := m.Bytes("text/html", content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func copyFile(sourceFileName string, destinationFileName string) error {
	if _, err := os.Stat(sourceFileName); err == nil {
		srcFile, err := os.Open(sourceFileName)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destinationFileName) // creates if file doesn't exist
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
		if err != nil {
			return err
		}

		err = destFile.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}

// CacheWrite stores a byte array in the expected destination
func CacheWrite(fileName string, content []byte, useGzip bool, minifyHTML bool, autoBackup bool) error {
	if autoBackup {
		err := copyFile(fileName, fileName+".old")
		if err != nil {
			return err
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	// var n int

	if minifyHTML {
		content, err = Minify(content)
		if err != nil {
			return err
		}
	}

	if useGzip {
		var b bytes.Buffer
		zippedContent := gzip.NewWriter(&b)
		zippedContent.Write(content)
		zippedContent.Close() // You must close this first to flush the bytes to the buffer.
		_, err = w.Write(b.Bytes())
		if err != nil {
			return err
		}
	} else {
		_, err = w.Write(content)
		if err != nil {
			return err
		}
	}
	w.Flush()
	// fmt.Printf("Caching content under [%s] (written bytes [%d])\n",fileName,n)
	return nil
}

// CacheRead loads cache content from file
func CacheRead(fileName string, useGzip bool) ([]byte, error) {
	var content []byte
	var err error
	if useGzip {
		content, err = ReadGzFile(fileName)
		if err != nil {
			return nil, err
		}
	} else {
		content, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}
