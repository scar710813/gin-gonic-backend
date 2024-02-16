#!/bin/bash

usage() {
    echo "Usage: $0 [-l|--local] [-p|--prod]" 1>&2
    exit 1
}

start_local() {
    docker-compose -f docker-compose-local.yml up
}

start_prod() {
    docker-compose up
}

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -l|--local) start_local;;
        -p|--prod) start_prod;;
        *) usage;;
    esac
    shift
done

