package main

import (
    "testing"
    "strings"
)

// Test whether "/etc" is a valid directory
// Should be true on any Unix system. Probably not for Windows
func Test_valid_dir__valid(t *testing.T) {
    valid := valid_dir("/etc")

    if !valid {
        t.Errorf("Valid directory is invalid")
    }
}

// Test whether "/abc" is a valid directory
// Should be false.
func Test_valid_dir__invalid(t *testing.T) {
    valid := valid_dir("/abc")

    if valid {
        t.Errorf("Invalid directory is valid")
    }
}

// Test whether "/usr/bin/whoami" is a valid file
// Should be true on any Unix system. Probably not for Windows
func Test_valid_file__valid(t *testing.T) {
    valid := valid_file("/usr/bin/whoami")

    if !valid {
        t.Errorf("Valid file is invalid")
    }
}

// Test whether "/usr/bin/xyzabc" is a valid file
// Should be false
func Test_valid_file__invalid(t *testing.T) {
    valid := valid_file("/usr/bin/xyzabc")

    if valid {
        t.Errorf("Invalid file is valid")
    }
}

// Test whether "/usr/bin/whoami" is a valid path
// Should be true
func Test_valid_path__file_valid(t *testing.T) {
    valid, _ := valid_path("/usr/bin/whoami")

    if !valid {
        t.Errorf("Valid file is invalid")
    }
}

// Test whether "/usr/bin/abcxyz" is a valid path
// Should be false
func Test_valid_path__file_invalid(t *testing.T) {
    valid, _ := valid_path("/usr/bin/abcxyz")

    if valid {
        t.Errorf("Invalid file is valid")
    }
}

// Test whether "/bin" is a valid path
// Should be true
func Test_valid_path__dir_valid(t *testing.T) {
    valid, _ := valid_path("/bin")

    if valid {
        t.Errorf("Invalid directory is valid")
    }
}

// Test whether "/abc" is a valid path
// Should be false
func Test_valid_path__dir_invalid(t *testing.T) {
    valid, _ := valid_path("/abc")

    if valid {
        t.Errorf("Invalid directory is valid")
    }
}

// Retrieve the file extension from "/etc/mime.types"
// Should be ".types"
func Test_get_extension__valid(t *testing.T) {
    ext := get_extension("/etc/mime.types")

    if ext != ".types" {
        t.Errorf("Valid extension is invalid")
    }
}

// Retrieve the file extension from "/bin"
// Should return empty string
func Test_get_extension__invalid(t *testing.T) {
    ext := get_extension("/bin")

    if ext != "" {
        t.Errorf("Invalid extension is valid")
    }
}

// Test whether 1023 bytes is formatted correctly
// Should return "1023 B"
func Test_format_bytes__bytes(t *testing.T) {
    t.Parallel()
    bytes := 1023
    formatted := format_bytes(int64(bytes))

    if formatted != "1023 B" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 2048 bytes is formatted correctly
// Should return "2.0 KB"
func Test_format_bytes__kilobytes(t *testing.T) {
    t.Parallel()
    bytes := 2048
    formatted := format_bytes(int64(bytes))

    if formatted != "2.0 KB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 1048576 bytes is formatted correctly
// Should return "1.0 MB"
func Test_format_bytes__megabytes(t *testing.T) {
    t.Parallel()
    bytes := 1048576
    formatted := format_bytes(int64(bytes))

    if formatted != "1.0 MB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 3221225472 bytes is formatted correctly
// Should return "3.0 GB"
func Test_format_bytes__gigabytes(t *testing.T) {
    t.Parallel()
    bytes := 3221225472
    formatted := format_bytes(int64(bytes))

    if formatted != "3.0 GB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 4398046511104 bytes is formatted correctly
// Should return "4.0 TB"
func Test_format_bytes__terabytes(t *testing.T) {
    t.Parallel()
    bytes := 4398046511104
    formatted := format_bytes(int64(bytes))

    if formatted != "4.0 TB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 3377699720527872 bytes is formatted correctly
// Should return "3.0 PB"
func Test_format_bytes__petabytes(t *testing.T) {
    t.Parallel()
    bytes := 3377699720527872
    formatted := format_bytes(int64(bytes))

    if formatted != "3.0 PB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 2305843009213693952 bytes is formatted correctly
// Should return "2.0 EB"
func Test_format_bytes__exabytes(t *testing.T) {
    t.Parallel()
    bytes := 2305843009213693952
    formatted := format_bytes(int64(bytes))

    if formatted != "2.0 EB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether get_mime_type returns correct MIME type from HTML snippet
// Should return string containing "text/html"
func Test_get_mime_type__html(t *testing.T) {
    mime := get_mime_type([]byte("<html><head><title>Hello</title></head><body><h1></h1></body></html>"))

    if !strings.Contains(mime, "text/html") {
        t.Errorf("MIME type produced: %s", mime)
    }
}

// Test whether get_mime_type returns correct MIME type from CSS snippet
// Should return string containing "text/plain"
func Test_get_mime_type__plain(t *testing.T) {
    mime := get_mime_type([]byte("body { font-size: 12px; }"))

    if !strings.Contains(mime, "text/plain") {
        t.Errorf("MIME type produced: %s", mime)
    }
}


// Test whether get_mime_type returns correct MIME type from XML snippet
// Should return string containind "text/xml" or "application/xml"
func Test_get_mime_type__xml(t *testing.T) {
    mime := get_mime_type([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))

    if !strings.Contains(mime, "text/xml") && !strings.Contains(mime, "application/xml") {
        t.Errorf("MIME type produced: %s", mime)
    }
}
