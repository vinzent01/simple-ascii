
echo "starting compilation"

echo "INFO: to compile to windows from linux needs mingw"
echo "INFO: to compile to mac from linux requires OSXCross"

echo "compiling linux..."
mkdir -p build
mkdir -p build/linux
cp assets build/linux/ -r
# Compilar para Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o ./build/linux/simple-ascii-linux ./src
7z a ./releases/linux ./build/linux


echo "compiling windows..."
mkdir -p build/windows
cp assets build/windows/ -r
# Compilar para Windows (x86_64)
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./build/windows/simple-ascii.exe ./src
7z a ./releases/windows ./build/windows



mkdir -p build/mac
cp assets build/mac/ -r
echo "compiling to mac..." 
CGO_ENABLED=1 CC=x86_64-apple-darwin21.1-clang GOOS=darwin GOARCH=amd64 go build -ldflags "-linkmode external -s -w '-extldflags=-mmacosx-version-min=10.15'" -o ./build/simple-ascii-mac ./src