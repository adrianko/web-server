package main

import (
    "testing"
    "github.com/adrianko/web-server"
)

func TestValid_dir(t * testing.T) {
    valid := valid_dir("/etc")

    if !valid {
        t.Errorf("valid_dir: /etc")
    }
}
