package html2pdf

/*
#cgo CFLAGS: -I./wkhtmltopdf/include
#include <wkhtmltox/pdf.h>
#include <stdio.h>

void progress_changed_cgo(wkhtmltopdf_converter *converter, int p) {
  progress_changed(converter, p);
}

void phase_changed_cgo(wkhtmltopdf_converter *converter) {
  phase_changed(converter);
}

void set_error_cgo(wkhtmltopdf_converter *converter, const char *msg) {
  set_error(converter, msg);
}

void set_warning_cgo(wkhtmltopdf_converter *converter, const char *msg) {
  set_warning(converter, msg);
}

*/
import "C"
