#! /bin/bash

if [[ $UID != 0 ]]; then
    echo "Please run with sudo:"
    echo "sudo $0 $*"
    exit 1
fi

service apache2 start
service nignx start
/home/vagrant/testing/./server &
