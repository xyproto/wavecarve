all: 
	@./build.sh
clean:
	@rm -f caire
install: all
	@cp caire /usr/local/bin
uninstall: 
	@rm -f /usr/local/bin/caire
package:
	@NOCOPY=1 ./build.sh package
test:
	go test -v -json ./... -run=. > ./test-report.json -coverprofile=coverage.out