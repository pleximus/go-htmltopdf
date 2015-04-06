# go-htmltopdf
Go bindings for wkhtmltopdf - Convert HTML to PDF using Webkit

## Usage

```go
package main

import "github.com/pleximus/go-htmltopdf"

converter := html2pdf.New()

// Converting HTML data and returning data buffer
// which can be used to write data to a file or 
// pipe it directly as http response
converter.SetData("<h2>Trying html to pdf</h2>")
err, data := converter.CreatePDF()
if err != nil {
  t.Fatal(err)
}
err = ioutil.WriteFile("try.pdf", data, 0x777)
if err != nil {
  t.Fatal(err)
}

// Fetching a URL and storing it in a PDF file
converter.SetURL("www.google.com")
converter.SetOutputFileName("google.pdf")
err, _ := converter.CreatePDF()
if err != nil {
  t.Fatal(err)
}
converter.Destroy()
```
