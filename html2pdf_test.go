package html2pdf

import "fmt"
import "testing"
import "io/ioutil"

func TestHtml2pdf(t *testing.T) {
	newObj := New()
	defer newObj.Destroy()

	newObj.SetData("Testing html to pdf go bindings")
	newObj.SetOutputFileName("test.pdf")

	newObj.OnProgressChanged(func(a int) { fmt.Println("Progress::", a, "%") })

	newObj.OnPhaseChanged(func(msg string) { fmt.Println("Phase::", msg) })

	newObj.OnError(func(msg string) { fmt.Println("Error::", msg) })

	newObj.OnWarning(func(msg string) { fmt.Println("Warning::", msg) })

	err, _ := newObj.CreatePDF()
	if err != nil {
		t.Fatal(err)
	}

	newObj.SetURL("www.google.com")
	newObj.SetBufferedOutput()
	err, data := newObj.CreatePDF()
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("google.pdf", data, 0x777)
	if err != nil {
		t.Fatal(err)
	}
}
