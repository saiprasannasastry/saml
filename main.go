package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://dev-6136203.okta.com/app/dev-6136203_test_4/exk3awcwssnpwhzc55d6/sso/saml", 301)
}

func blah(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "https://dev-6136203.okta.com/app/dev-6136203_test_4/exk3awcwssnpwhzc55d6/sso/saml", 301)
	r.ParseForm()

	for k, v := range r.Form {
		if k == "SAMLResponse" {
			dc, err := b64.StdEncoding.DecodeString(v[0])
			if err != nil {
				fmt.Printf("The err is %v", err)
			}
			resp := Decoder{}
			err = xml.Unmarshal(dc, &resp)
			//fmt.Printf("The resp is %+v,", resp.Assertion.Signature.KeyInfo.Cert.X509Certificate)
			file, err := ioutil.TempFile("/Users/saiprasanna.sastryss/go/src/github1/saml", "*.cert")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(file.Name())
			begin := []byte("-----BEGIN CERTIFICATE-----\n")
			end := []byte("\n-----END CERTIFICATE-----\n")
			cert := []byte(resp.Assertion.Signature.KeyInfo.Cert.X509Certificate)

			if _, err = file.Write(begin); err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}

			if _, err = file.WriteString(insertNth(string(cert), 76)); err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}
			if _, err = file.Write(end); err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}
			newfile, err := os.Open(file.Name())
			if err != nil {
				log.Fatal(err)
			}

			b, _ := ioutil.ReadAll(newfile)

			inputCert, err := os.Open("azure.cert")
			if err != nil {
				log.Fatal(err)
			}
			c, _ := ioutil.ReadAll(inputCert)

			same := strings.Compare(strings.TrimSpace(string(b)), strings.TrimSpace(string(c)))
			var roles []AttributeValue
			var logout string
			if same == 0 {
				logout = resp.Assertion.Conditions.NotOnOrAfter
				for _, v := range resp.Assertion.AttributeStatement.Attribute {
					if v.Name == "user.rolew" {
						roles = v.AttributeValue
					}
				}
			}
			fmt.Printf("The roles are %v, logout value is %v", roles, logout)

		}
	}

	fmt.Fprintf(w, "Hello, World!")
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}

func main() {
	app := http.HandlerFunc(blah)
	http.Handle("/hello", app)
	http.HandleFunc("/hello/test", hello)
	http.ListenAndServe(":8081", nil)
	fmt.Println("hello")
}
