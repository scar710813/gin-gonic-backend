#!/bin/bash

usage() {
    echo "Usage: $0 [-p|--push] [-b|--build]" 1>&2
    exit 1
}

build_image() {
    docker build -t whyakari/shinsei:latest .
}

push_image() {
    docker push whyakari/shinsei:latest
}

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -b|--build) build_image;;
        -p|--push) push_image;;
        *) usage;;
    esac
    shift
done

