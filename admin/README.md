# Regenerating the Admin service

## Install
- goagen: `go install github.com/goadesign/goa/goagen`
- go-bindata: `go get -u github.com/jteeuwen/go-bindata/...`

## Run
```
goagen controller --regen --pkg admin -d github.com/zenoss/zenkit/admin/design
goagen app -d github.com/zenoss/zenkit/admin/design
goagen swagger -d github.com/zenoss/zenkit/admin/design
go-bindata -ignore='swagger\.go' -pkg swagger -o swagger/swagger.go swagger/
```