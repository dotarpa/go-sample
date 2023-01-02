// https://github.com/prometheus/alertmanager/blob/main/notify/email/email.go
// Just looking into mime/quotedprintable.

// The Go Playground
// https://go.dev/play/p/W7a7tBaRNIU

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

func main() {
	multipartBuffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(multipartBuffer)
	boundary := multipartWriter.Boundary()

	w, err := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Type":              {"text/plain; charset=UTF-8"},
	})

	if err != nil {
		log.Fatal(err)
	}

	body := "Hello, 世界: RFC2822で規定されている78文字以上のメッセージをうまいこと分割するmime/quotedprintableの使い方を調べてるだけ"

	qw := quotedprintable.NewWriter(w)
	_, err = qw.Write([]byte(body))

	if err != nil {
		log.Fatal(err)
	}

	err = qw.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = multipartWriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("boundary: %v\n", boundary)
	fmt.Println(multipartBuffer.String())

	dec_body := strings.NewReader(multipartBuffer.String())
	multipartReader := multipart.NewReader(dec_body, boundary)
	for {
		p, err := multipartReader.NextPart()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		slurp, err := io.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(slurp))
	}
}
