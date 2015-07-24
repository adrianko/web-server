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
// SHould return empty string
func Test_get_extension__invalid(t *testing.T) {
    ext := get_extension("/bin")

    if ext != "" {
        t.Errorf("Invalid extension is valid")
    }
}
