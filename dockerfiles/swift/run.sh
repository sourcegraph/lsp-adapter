#!/bin/bash

# until we have linux support, this script automates building and running on macOS
# https://github.com/RLovelett/langserver-swift/issues/18

set -eax

if [ ! -d langserver-swift ]; then
    git clone https://github.com/RLovelett/langserver-swift.git
fi

COMMIT=838d8fadfc9fad82c8c1417e2f70a7d73f6347ee
cd langserver-swift
git checkout $COMMIT || (git fetch && git checkout $COMMIT)
make debug

lsp-adapter -trace -didOpenLanguage=swift .build/x86_64-apple-macosx10.10/debug/langserver-swift
