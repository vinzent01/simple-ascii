
echo "starting compilation"

echo "INFO: to compile to windows from linux needs mingw"
echo "INFO: to compile to mac from linux requires OSXCross"

echo "compiling linux..."
# Compilar para Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o ./build/simple-ascii-linux ./src

echo "compiling windows..."
# Compilar para Windows (x86_64)
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./build/simple-ascii.exe ./src

echo "compiling to mac..." 
CGO_ENABLED=1 CC=x86_64-apple-darwin21.1-clang GOOS=darwin GOARCH=amd64 go build -ldflags "-linkmode external -s -w '-extldflags=-mmacosx-version-min=10.15'" -o ./build/simple-ascii-mac ./src