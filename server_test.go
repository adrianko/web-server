package main

import (
    "testing"
)

func Test_valid_dir(t * testing.T) {
    valid := valid_dir("/etc")

    if !valid {
        t.Errorf("valid_dir: /etc")
    }
}
