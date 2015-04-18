package html2pdf

/*
#cgo CFLAGS: -I./wkhtmltopdf/include
#cgo LDFLAGS: -lwkhtmltox -L./wkhtmltopdf/lib
#include <wkhtmltox/pdf.h>

void progress_changed_cgo(wkhtmltopdf_converter *converter, int p);
void phase_changed_cgo(wkhtmltopdf_converter *converter);
void set_error_cgo(wkhtmltopdf_converter *converter, const char *msg);
void set_warning_cgo(wkhtmltopdf_converter *converter, const char *msg);

*/
import "C"
import "errors"
import "unsafe"

var inputData string

type intFuncCallback func(int)
type stringFuncCallback func(string)

type Html2pdf struct {
	gs        *C.wkhtmltopdf_global_settings
	os        *C.wkhtmltopdf_object_settings
	converter *C.wkhtmltopdf_converter
}

var progressChangedCallback intFuncCallback
var phaseChangedCallback, errorCallback, warningCallback stringFuncCallback

//export progress_changed
func progress_changed(converter *C.wkhtmltopdf_converter, p C.int) {
	if progressChangedCallback != nil {
		progressChangedCallback(int(p))
	}
}

//export phase_changed
func phase_changed(converter *C.wkhtmltopdf_converter) {
	var phase C.int = C.wkhtmltopdf_current_phase(converter)
	var phaseDescription = C.GoString(C.wkhtmltopdf_phase_description(converter, phase))
	if phaseChangedCallback != nil {
		phaseChangedCallback(phaseDescription)
	}
}

//export set_error
func set_error(converter *C.wkhtmltopdf_converter, msg *C.char) {
	if errorCallback != nil {
		errorCallback(C.GoString(msg))
	}
}

//export set_warning
func set_warning(converter *C.wkhtmltopdf_converter, msg *C.char) {
	if warningCallback != nil {
		warningCallback(C.GoString(msg))
	}
}

func setProgressChangedCallback(h2p *Html2pdf) {
	C.wkhtmltopdf_set_progress_changed_callback(h2p.converter, C.wkhtmltopdf_int_callback(unsafe.Pointer(C.progress_changed_cgo)))
}

func setPhaseChangedCallback(h2p *Html2pdf) {
	C.wkhtmltopdf_set_phase_changed_callback(h2p.converter, C.wkhtmltopdf_void_callback(unsafe.Pointer(C.phase_changed_cgo)))
}

func setErrorCallback(h2p *Html2pdf) {
	C.wkhtmltopdf_set_error_callback(h2p.converter, C.wkhtmltopdf_str_callback(unsafe.Pointer(C.set_error_cgo)))
}

func setWarningCallback(h2p *Html2pdf) {
	C.wkhtmltopdf_set_error_callback(h2p.converter, C.wkhtmltopdf_str_callback(unsafe.Pointer(C.set_warning_cgo)))
}

func (h2p *Html2pdf) OnProgressChanged(callback intFuncCallback) {
	progressChangedCallback = callback
}

func (h2p *Html2pdf) OnPhaseChanged(callback stringFuncCallback) {
	phaseChangedCallback = callback
}

func (h2p *Html2pdf) OnError(callback stringFuncCallback) {
	errorCallback = callback
}

func (h2p *Html2pdf) OnWarning(callback stringFuncCallback) {
	warningCallback = callback
}

func New() *Html2pdf {
	C.wkhtmltopdf_init(1)
	gs := C.wkhtmltopdf_create_global_settings()
	os := C.wkhtmltopdf_create_object_settings()
	return &Html2pdf{gs: gs, os: os}
}

func createConverter(h2p *Html2pdf) {
	converter := C.wkhtmltopdf_create_converter(h2p.gs)
	h2p.converter = converter
	setProgressChangedCallback(h2p)
	setPhaseChangedCallback(h2p)
	setErrorCallback(h2p)
	setWarningCallback(h2p)
}

func (h2p *Html2pdf) SetGlobalSettings(global_settings [][2]string) {
	for _, row := range global_settings {
		C.wkhtmltopdf_set_global_setting(h2p.gs, C.CString(row[0]), C.CString(row[1]))
	}
}

func (h2p *Html2pdf) SetObjectSettings(object_settings [][2]string) {
	for _, row := range object_settings {
		C.wkhtmltopdf_set_object_setting(h2p.os, C.CString(row[0]), C.CString(row[1]))
	}
}

func (h2p *Html2pdf) SetURL(url string) {
	C.wkhtmltopdf_set_object_setting(h2p.os, C.CString("page"), C.CString(url))
}

func (h2p *Html2pdf) SetOutputFileName(name string) {
	C.wkhtmltopdf_set_global_setting(h2p.gs, C.CString("out"), C.CString(name))
}

func (h2p *Html2pdf) SetBufferedOutput() {
	C.wkhtmltopdf_set_global_setting(h2p.gs, C.CString("out"), C.CString(""))
}

func (h2p *Html2pdf) SetData(data string) {
	inputData = data
}

func (h2p *Html2pdf) CreatePDF() (error, []byte) {
	createConverter(h2p)

	if inputData == "" {
		C.wkhtmltopdf_add_object(h2p.converter, h2p.os, nil)
	} else {
		C.wkhtmltopdf_add_object(h2p.converter, h2p.os, C.CString(inputData))
		inputData = ""
	}

	res := C.wkhtmltopdf_convert(h2p.converter)

	if res != 1 {
		return errors.New("Conversion failed!"), nil
	}

	var ptr *C.uchar
	length := C.wkhtmltopdf_get_output(h2p.converter, &ptr)

	// GoBytes accepts C.int as length, so the length needs to be
	// type casted from long to int
	outData := C.GoBytes(unsafe.Pointer(ptr), C.int(length))
	C.wkhtmltopdf_destroy_converter(h2p.converter)
	return nil, outData
}

func (h2p *Html2pdf) Destroy() {
	C.wkhtmltopdf_deinit()
}
