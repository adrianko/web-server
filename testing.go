package test

import (
    "testing"
    "github.com/adrianko/web-server/src"
)

func TestValid_dir(t * testing.T) {
    valid := valid_dir("/etc")

    if !valid {
        t.Errorf("valid_dir: /etc")
    }
}
