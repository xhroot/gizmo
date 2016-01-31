#!/bin/bash

set -e

workdir=.cover
profile="$workdir/cover.out"
mode=set

generate_cover_data() {
    rm -rf "$workdir"
    mkdir "$workdir"

    for pkg in "$@"; do
        if [ $pkg == "github.com/xhroot/gizmo/server" -o $pkg == "github.com/xhroot/gizmo/config" -o $pkg == "github.com/xhroot/gizmo/web" -o $pkg == "github.com/xhroot/gizmo/pubsub" ]
            then
                f="$workdir/$(echo $pkg | tr / -)"
                go test -covermode="$mode" -coverprofile="$f.cover" "$pkg"
        fi
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

show_cover_report() {
    go tool cover -${1}="$profile"
}

push_to_coveralls() {
    goveralls -coverprofile="$profile"
}

generate_cover_data $(go list ./...)
show_cover_report func
case "$1" in
"")
    ;;
--html)
    show_cover_report html ;;
--coveralls)
    push_to_coveralls ;;
*)
    echo >&2 "error: invalid option: $1" ;;
esac
rm -rf "$workdir"
