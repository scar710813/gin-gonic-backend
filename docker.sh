#!/bin/bash

usage() {
    echo "Usage: $0 [-p|--push] [-b|--build]" 1>&2
    exit 1
}

build_image() {
    docker build -t whyakari/rinha:2.0 .
}

push_image() {
    docker push whyakari/rinha:2.0
}

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -b|--build) build_image;;
        -p|--push) push_image;;
        *) usage;;
    esac
    shift
done

