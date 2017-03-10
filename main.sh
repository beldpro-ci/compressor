#!/bin/bash

set -o errexit

main () {
  case $1 in
    setup)
      test "$TEST_SUITE" == "integration" && install_docker
      install_go_dependencies
      make deps
      ;;

    test)
      test "$TEST_SUITE" == "integration" \
        && make test -j8 \
        || make test-unit -j2
      ;;

    *)
      echo "Unknown option '$1'"
      exit 1
      ;;
  esac
}

install_go_dependencies () {
  go get -u -v github.com/Masterminds/glide
  go get -u -v github.com/jteeuwen/go-bindata/...
  go get -u -v github.com/hashicorp/consul
}


install_docker () {
  sudo apt-get update -y
  sudo apt-get install --only-upgrade docker-engine -y
}


main "$@"

