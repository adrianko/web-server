package main

import (
    "testing"
)

// Test whether "/etc" is a valid directory
// Should be true on any Unix system. Probably not for Windows
func Test_valid_dir_valid(t *testing.T) {
    valid := valid_dir("/etc")

    if !valid {
        t.Errorf("Valid directory is invalid")
    }
}

// Test whether "/abc" is a valid directory
// Should be false.
func Test_valid_dir_invalid(t *testing.T) {
    valid := valid_dir("/abc")

    if valid {
        t.Errorf("Invalid directory is valid")
    }
}

// Test whether "/usr/bin/whoami" is a valid file
// Should be true on any Unix system. Probably not for Windows
func Test_valid_file_valid(t *testing.T) {
    valid := valid_file("/usr/bin/whoami")

    if !valid {
        t.Errorf("Valid file is invalid")
    }
}

// Test whether "/usr/bin/xyzabc" is a valid file
// Should be false
func Test_valid_file_invalid(t *testing.T) {
    valid := valid_file("/usr/bin/xyzabc")

    if valid {
        t.Errorf("Invalid file is valid")
    }
}
