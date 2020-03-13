moduleowner = github.com/lmnzr/
packagename = mygosql
# Build And Development
init:
	@ go mod init $(moduleowner)$(packagename)
	@ go mod vendor 
update:
	@ go mod vendor
clean:
	@ sudo rm -rf $(packagename).bin $(packagename).exe cover.txt cover.html cover.out build
test:
	@ go test $(moduleowner)$(packagename)/test/... 
test-cover:
	@ mkdir -p build
	@ go test $(moduleowner)$(packagename)/test/... -coverpkg=./... -coverprofile=./build/cover.out
	@ go tool cover -html=./build/cover.out -o ./build/cover.html   
run:
	@ go build -o ./build/$(packagename).bin $(moduleowner)$(packagename)/cmd/$(packagename)  && ./build/$(packagename).bin	  
swagger:
	@ cd cmd/$(packagename) && swag init
.PHONY: init update run test test-cover clean swagger