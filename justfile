build-copy:
    mkdir -p dist
    go build -o dist/tscp ./cmd/tscp
    sudo cp dist/tscp /usr/local/bin

dry-run:
    cd test && go run ../cmd/tscp --ext ".flac" --dry-run test
