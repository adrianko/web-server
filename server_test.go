package main

import "testing"

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
    bytes := 1023
    formatted := format_bytes(int64(bytes))

    if formatted != "1023 B" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 2048 bytes is formatted correctly
// Should return "2.0 KB"
func Test_format_bytes_kilobytes(t *testing.T) {
    bytes := 2048
    formatted := format_bytes(int64(bytes))

    if formatted != "2.0 KB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 1048576 bytes is formatted correctly
// Should return "1.0 MB"
func Test_format_bytes_megabytes(t *testing.T) {
    bytes := 1048576
    formatted := format_bytes(int64(bytes))

    if formatted != "1.0 MB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}

// Test whether 3221225472 bytes is formatted correctly
// Should return "3.0 GB"
func Test_format_bytes_gigabytes(t *testing.T) {
    bytes := 3221225472
    formatted := format_bytes(int64(bytes))

    if formatted != "3.0 GB" {
        t.Errorf("Format bytes produced: %s from %d", formatted, bytes)
    }
}
