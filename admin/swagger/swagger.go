// Code generated by go-bindata.
// sources:
// swagger/swagger.json
// swagger/swagger.yaml
// DO NOT EDIT!

package swagger

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _swaggerSwaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x57\xdd\x6e\x1b\x37\x13\x7d\x95\x01\xfd\x5d\x7c\x05\x36\xb6\xf2\x8b\x56\x45\x2f\x0a\xb8\x41\xdc\xb4\x4d\x60\x27\xb9\x09\x82\x60\xb4\x1c\x49\x0c\xb8\x24\x4d\x0e\x65\x0b\x81\xde\xbd\x18\xee\xae\xac\xb5\xa4\x58\x69\x5d\x24\x57\x6b\xef\x70\x87\xe7\x9c\x39\x33\xa4\x3e\xab\x74\x85\xb3\x19\x45\x35\x56\x8f\x8e\x47\xaa\x52\xc6\x4d\xbd\x1a\x7f\x56\x6c\xd8\x92\x1a\xab\x5f\x75\x63\x1c\x5c\x50\x5c\x98\x9a\x54\xa5\x34\xa5\x3a\x9a\xc0\xc6\x3b\x35\x56\x6f\xd9\x58\xc3\x86\x12\x84\xe8\x17\x46\x93\x86\xc9\x12\x78\x4e\x80\xe5\x3b\x72\x3a\x78\xe3\x58\x55\x6a\x41\x31\xb5\x1f\xa9\x55\xa5\x52\x3d\xa7\x86\x92\x1a\xbf\x57\x73\xe6\xa0\x3e\x54\xaa\xf6\x2e\xe5\xee\x1d\x86\x60\x4d\x8d\xb2\xcb\xc9\xa7\xe4\x9d\xc4\x43\xf4\x3a\xd7\x5f\x88\x23\xcf\x93\x40\x3f\x99\x13\x5a\x9e\xcb\x9f\x33\xe2\x42\x06\x67\xed\x56\x6d\xe0\x43\xa5\x52\x6e\x1a\x8c\x4b\x35\xee\xde\x41\x17\xba\x4d\xf0\x9c\x82\x8f\x5c\x18\x75\x0b\xfd\xb4\xfc\x97\xd6\x8a\xf8\x40\xb1\x20\x39\xd3\xeb\x74\x47\xdd\xe3\xa1\x1a\xe2\x66\xba\xe6\x93\x60\xd1\x14\xc4\x91\x52\xf0\x2e\x51\x41\xfd\x68\x34\x92\xc7\x70\xfb\x57\x2f\x45\xac\xa7\xa3\xc7\xdb\xa1\xae\x24\xf0\xd6\xe1\x02\x8d\xc5\x89\x25\xb5\xda\xa5\xec\xaa\x12\x50\xfa\x5b\xea\xf0\x3d\xa8\xb0\xaa\x7a\x63\x9c\x68\x7f\xe5\x24\x55\xf0\xe9\x4e\x7b\xc8\xda\x7d\xa2\x5c\x10\x27\x68\xd0\x65\xb4\x1f\x65\x97\x8f\x89\x91\x73\x02\xf6\x80\x0e\x28\x46\x1f\xf7\xe9\x52\x20\x7c\x49\x95\x80\x11\x1b\x62\x8a\x12\xfb\xac\x1c\x36\xd2\x8d\x01\x97\xd6\xa3\x2e\x7d\xaa\xc6\x6a\xe2\xf5\x52\x89\x82\x97\xd9\x44\xd2\x6a\xcc\x31\x53\x47\x1e\x85\xd9\xff\x22\x4d\xd5\x58\x1d\x9d\x68\x9a\x1a\x67\x04\x45\x3a\x39\xf5\x57\xee\x45\x81\xf1\xba\x4b\xb7\x5a\x1d\x5c\x87\xbb\xb4\xcd\xe1\x50\x65\x73\xf8\x07\xba\x3a\x63\xf7\x49\x9a\xc3\x3d\xd8\x6c\x2f\xbd\x86\x38\x9a\x3a\xed\x18\x2a\x65\xd2\x0d\xa9\x75\x8b\xdb\x21\xb8\xa3\x97\x38\x47\x07\x08\xc9\x61\x48\x73\xcf\xd2\x4b\x7d\xfe\xdb\xe4\x4a\x8a\xa3\x9b\xe8\x17\xa7\x60\x35\x78\xb5\x70\xfa\x78\xe6\xf1\xb8\x35\xe2\x7e\x4f\x45\x62\x5e\xf6\x96\xba\xcc\x14\x97\x5b\x90\xcf\x9c\x26\xc7\x10\x29\x65\xcb\xc6\xcd\xe0\xf7\x8b\x57\x7f\x0d\x9c\x37\x45\x9b\xa8\x52\xbc\x0c\x54\x8c\xe9\x2d\x61\x4b\x7d\x8a\xd9\x72\x6b\xcd\xc3\x6d\x26\xed\xbe\x23\x74\xe6\x98\xa2\x43\x5b\x0e\x24\x8a\xf0\x5b\xd7\x63\x77\x38\xbe\x55\x60\xb5\xbf\xba\xc1\xb8\xd9\x61\xa5\x95\x95\x7b\xeb\x2a\xcc\x34\x5c\x19\x9e\x03\xc2\xa3\xd1\x08\xcc\x60\x48\x82\x49\x70\x33\xa6\x76\x57\x5a\xf2\xdf\xcb\xb1\x71\xe0\x51\xf0\x8d\x28\xfe\x97\xad\xba\xbe\xd2\x6c\xd5\xb3\x8f\x0c\xe8\x76\x2f\xa1\x0f\xde\xe6\x7c\x6a\x52\xb0\xb8\x84\x8b\x6e\x5d\x4e\xa2\xcf\x39\x9d\xfa\x7a\x8b\x60\x97\xe3\xe8\x26\xd7\x16\xcb\x39\x37\xf6\xfe\x48\x1e\x97\xbe\x3f\x90\xa9\xac\xdd\x4b\xf3\x5c\x46\x0c\x2d\x68\xcd\x33\x05\xaa\x01\x53\xdf\xea\xbb\x89\x76\x63\xe7\xdf\x4c\xa5\xef\x6f\x20\xac\xca\xd8\xea\x17\x4b\x92\xed\x23\x73\xe3\x8a\xbc\x1d\x5c\xcf\x41\x3f\xf9\x44\x35\xb7\x02\x05\x8a\x72\x53\x96\x2f\x23\x61\x57\xb6\x6e\x5d\xe2\xd8\x36\x05\x5d\x63\x13\x4a\xd6\x3f\x30\x6b\x74\x6c\x72\x03\x97\xd9\x00\x5d\xcb\x63\x92\x93\xc6\x06\x12\x06\x43\x8e\x09\x98\x9a\xe0\x23\x1e\x17\xa3\xac\xbf\xbd\xd9\xe0\x6b\xb3\x6c\x4e\xf4\xf7\x7d\x16\x19\x1b\xad\x60\x1b\xa4\xff\x24\x6d\x50\xd0\x83\x91\xa3\xc1\x4c\x0d\xc5\x31\xec\x2d\xf4\xcf\xb0\x30\x74\xf5\x4b\x7f\x1a\xdc\x25\x50\xed\x35\xed\x92\x67\x58\x7c\x74\x9b\x1b\x3e\x10\xc3\x9a\xa9\xa9\xdb\x9b\x17\x48\x8e\x0a\xe8\x3a\x44\x4a\x89\xb4\x58\x19\xa1\xcd\x04\x0b\xb4\x99\x8e\x07\x7a\x1b\xb7\x40\x6b\xf4\xc7\x12\x52\xc5\x02\x8c\xc6\x1e\x80\x02\xe6\xb9\x41\xf7\x20\x12\x6a\x99\x7a\xb2\xa5\x45\x57\x30\xc1\x1a\x13\x7b\xe0\xb9\x49\xe0\xeb\x3a\xc7\x48\xae\xa6\xfe\x06\x1d\xa2\x9f\x58\x6a\x86\x68\xde\x09\x0a\x59\x71\x76\x0a\x4d\x4e\x0c\x13\x92\x3b\xa5\x71\x4c\xd2\xbc\xab\x4a\x19\x7d\x08\xb4\xec\xcc\x65\xde\xac\x11\x4c\x7d\x6c\x91\x04\x8c\x6c\xea\x6c\x31\x1e\x0a\xea\xf1\xf3\x87\xcf\x5f\xbe\x3b\x3f\x97\xed\x1b\x62\xdc\x00\xb0\xae\xe3\x6d\x00\xb2\x0e\xda\x28\xd4\xde\x31\x1a\x27\x05\x70\x52\x2f\x46\xa7\x31\xea\xb2\xe6\x81\xfc\xf0\x8c\x4d\xab\x1a\x4e\x7c\x6e\x7f\x6c\x94\x52\x0e\x50\x88\x05\x1b\x4a\x8c\x4d\x50\xe3\x87\x4f\x9e\xfe\xf8\x6c\xf4\xd3\xe8\xd9\xb3\x55\xa5\x50\xeb\xd2\xb1\x68\x5f\x6f\x98\xa9\xdc\x3b\x2a\xd5\x5e\x21\xef\x96\x4c\x36\x7d\xf1\xe6\xcd\x6b\xe8\xee\x9c\xe2\xa2\xde\x65\x52\xdc\xbe\x8e\x9d\x40\x5f\x61\xb0\x27\xa3\x91\x6a\x47\xcb\xe6\x7e\x65\x60\x41\x3f\x03\xa1\x91\xb6\x82\xd2\x57\xff\xef\x9a\xa5\x74\xce\x0f\x43\x09\xda\xfe\xb8\xe5\xd9\x1b\xcb\xde\x65\x9f\xd6\x3d\x37\xf5\xdc\x28\xe7\x1e\x6d\x7b\xf9\x3a\x16\xab\x5b\x63\xfb\xd5\xcb\x3d\xc7\xd7\xea\xef\x00\x00\x00\xff\xff\xcc\x95\x1b\x1f\x62\x10\x00\x00")

func swaggerSwaggerJsonBytes() ([]byte, error) {
	return bindataRead(
		_swaggerSwaggerJson,
		"swagger/swagger.json",
	)
}

func swaggerSwaggerJson() (*asset, error) {
	bytes, err := swaggerSwaggerJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger/swagger.json", size: 4194, mode: os.FileMode(420), modTime: time.Unix(1507300625, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _swaggerSwaggerYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x98\x4f\x6f\x1b\x37\x13\xc6\xef\xfa\x14\x03\xf9\x05\xfc\x16\xad\x25\xe5\x5f\xd1\x6c\xd1\x43\x01\x37\x88\xeb\xb6\x0e\xec\x24\xd7\x60\xb4\x1c\x49\x13\x70\x49\x9a\x1c\xca\xd6\xb7\x2f\x48\x51\xd2\xae\xb5\x56\x9c\x20\x35\x5a\xb4\xba\x2c\x76\x39\x4b\x3e\xfc\xcd\xc3\x21\x57\xb5\x35\x21\x36\x14\xaa\xc1\x09\xa0\x73\x9a\x6b\x14\xb6\x66\xfc\x31\x58\x33\x50\x34\x63\xc3\xe9\x3e\x54\x03\x80\x53\x7b\x63\x5e\x13\x6a\x59\xbc\xc1\x95\xb6\xa8\xd2\x43\x00\xba\xc5\xc6\x69\x5a\xdf\x00\x78\xc2\x60\x4d\x05\xbf\x61\x54\x68\x84\x63\x03\xd7\x91\x81\x6e\xd3\x65\x1a\x83\xc2\x06\x02\x3a\x26\x23\x04\x42\x8d\xb3\x1e\x47\xf9\x5d\xe7\xad\x23\x2f\x4c\xe1\x4e\x5f\xe5\x6e\x37\xd4\x17\x74\x9e\x7e\xb2\x72\x54\x41\x10\xcf\x66\x3e\x58\xf7\x7f\x1d\xd9\x53\x99\xc8\x49\x19\x2f\xdf\x08\x4b\x1a\x68\x6f\xce\x83\x5d\x47\x76\xfa\x91\x6a\x19\x00\x90\xf7\xd6\xaf\xfb\x50\x14\x6a\xcf\x2e\x31\xab\xe0\x97\xf4\x1c\x3c\x05\x67\x4d\x20\x68\x48\x31\xe6\x77\xe1\xff\x8a\x66\x18\xb5\xc0\x92\xe9\xe6\x9b\x3e\x8c\xb5\x55\x54\x01\x9b\x25\x6a\x56\x1f\x96\xa8\x23\x95\x16\x45\x82\xac\x2b\x78\x9f\x9e\x81\x9d\xc1\xd9\x29\x34\x31\x08\x4c\x09\xd0\x00\x1b\xa1\x39\xf9\x12\xcc\xaa\x82\x67\xaf\x9e\xbc\x3a\x7f\x7f\x79\x59\x1e\x35\x24\xb8\x43\x2a\xdc\x50\x10\x6c\x5c\x05\x4f\x46\xcf\x5f\xfc\xf0\xfd\xe4\x25\x7d\x3b\x79\x59\xda\x83\xa0\xc4\x50\xc1\xf0\xf9\x64\x32\xbc\x27\x49\x59\xe9\xb6\xbf\x0e\x00\x34\x6d\x53\x9d\x04\x47\x35\xcf\xb8\x5e\x03\xcb\x2f\x7e\x07\x74\xeb\x3c\x85\x40\x0a\x30\x00\x96\xec\x40\x9e\xf0\x68\x3f\xf1\x7d\x40\x7a\x12\xbb\xa5\x74\x8f\x2e\x58\xc4\x06\xcd\x89\x27\x54\x38\xd5\x94\x44\x68\x34\x59\x25\x6c\x55\x8a\x05\x59\x70\x00\x5b\xd7\xd1\x7b\x32\x75\xa6\x2d\x8b\xdd\xb8\x19\xc7\x54\x53\xd3\xa3\xf4\x41\xe9\xe9\x95\xce\xea\x5e\xd9\xd1\xf0\x75\x24\x60\x45\x46\x78\xc6\xe4\x61\x66\xfd\x5a\xa5\x43\x2f\x5c\x47\x8d\x7e\x5f\xf0\x01\x99\x77\xcc\xd1\x2b\xa8\xeb\x18\x54\x2a\xd7\x04\xd4\x6f\x76\x5e\x00\xf1\xad\x7c\xdc\x51\x9d\xde\x2f\xab\x05\x6a\x6b\x04\xd9\xa4\x1c\x9b\x64\x09\x41\xa3\xd0\xab\x1c\x73\xc2\x66\x66\x7d\xb3\x4e\x03\x4e\x6d\x94\x16\xe9\x34\x8d\xec\x9b\xfd\x49\xb4\xa3\x0e\xda\x79\x6f\xe5\xb6\x3c\xde\x2f\x3e\x8d\xfa\xfa\xed\xdb\x37\x25\x2a\x7b\x76\xe3\xe9\x64\x9c\x8d\x47\x0a\xe0\xae\x9d\x5b\xba\x3e\x69\xec\xdd\x0a\xeb\x4d\x42\x29\x47\xc7\xbf\xa7\x1a\x92\x4b\xc8\xce\x03\x55\xa7\x72\x2f\x8d\x1a\xcd\x2d\x8e\x32\xab\x1f\x73\x81\xf9\xa9\x54\x9b\xe3\xfd\xe2\x95\x88\xa7\xa9\x77\x26\xfd\x4e\x58\x73\x4a\x6b\x9a\xd6\x92\x15\x29\x98\xae\x32\x0a\x54\x0d\x1b\x20\xa3\x9c\x65\x93\x00\x16\x5d\x3f\xe7\xe7\x57\xe4\x97\x5c\x27\x1b\x2c\xc9\x87\xdc\xd5\x70\x38\x70\x28\x8b\x8c\x77\xbc\xc8\x85\x74\x4d\x7a\x4e\x52\x0d\x7a\x80\x5f\x92\xb3\x5e\xf2\x60\xeb\xf0\x8d\x85\xc3\xb6\xf3\xf4\x4b\xc6\xcb\xf3\x3d\x53\x55\x09\x3c\x2a\x97\x27\x83\xed\xd2\x54\xb1\xde\xd5\xa9\x13\x10\xba\x95\xb1\xd3\xc8\x66\xbb\xbf\xac\x4b\x73\x2b\xfb\xc3\xa7\x93\xc9\xb0\xed\xa8\x8e\xba\x8b\xf3\x5d\xe0\x8b\xc9\xb3\xfb\x03\x0b\x0a\x78\x67\x70\x89\xac\x93\x57\x36\x76\xab\x17\xd4\xb4\x55\x2d\x44\xdc\xa6\x2d\x36\x0d\xfa\xd5\x66\x46\xe5\x52\x1a\x05\xe7\xed\xb7\x76\x4d\x0b\x42\xf5\xd7\xb0\xfc\x77\x91\x2c\xfe\x1c\x2b\x7b\x53\xce\x1c\xce\x86\x7e\x97\x5e\x91\x04\x68\xd0\x44\xd4\x1f\xd2\xa8\x1f\x4a\x79\x10\x9b\x6a\x7c\x5e\x7a\x07\xe0\xa6\x01\x36\x68\xd1\x63\x43\x42\xbe\x25\x88\x4d\x05\x53\xab\x56\x5b\x22\x06\x1b\xaa\xc0\xb5\x0e\x20\xd0\x3e\xbc\x74\x6b\x6f\xa6\x82\x6d\x9c\xff\xf3\x34\xab\xe0\xf8\x68\xdc\x3a\xcf\x8d\xf7\x0e\x36\xc7\x8f\x91\xec\x07\xa5\x2c\xe1\xf9\xac\x84\x45\xf7\xe5\xe9\x32\xac\x0f\x64\x2a\xba\xbf\x0d\x95\xe8\x1e\xc6\xa4\x21\xf1\x5c\x87\x4f\x55\x59\x89\xde\xa4\x5d\xc9\xa0\x0b\x0b\x2b\xa9\x32\x94\x37\xfb\x70\xe4\xba\x7f\xd4\x0d\xe8\xb3\x6e\xd9\x68\x0e\x1d\x07\xce\x4c\xda\xb8\x12\xad\xa8\x25\xed\x88\xbf\x5e\x5d\xfc\xb1\x0d\x4d\xde\xbf\x8e\xe4\xf7\xcc\xef\x49\x64\xd5\xe3\xfd\x19\xea\x70\xf7\x20\x38\xb5\x56\x13\x9a\x7b\x53\xb7\xf7\xa1\xd3\xd7\xd0\xd9\x47\xbf\x4e\x95\x3b\x10\x78\x66\x84\xbc\x41\x9d\xcb\x1d\xf9\xf5\xa7\x43\x2b\x78\x7f\x55\xdf\xb3\xae\xb3\xdc\xe3\xcf\x31\x56\xc9\xea\x3a\xc7\xbd\xe6\xda\xb4\x8c\x1d\x9b\xf9\xa7\x8c\x95\x08\x29\xb8\x61\x59\x00\xc2\xd3\xc9\x04\xb8\xb3\xe7\x00\x07\xb8\x5b\xc8\x7b\xac\x96\x46\x7a\x9c\x9d\xfc\x41\x8c\x92\x9c\x07\x00\x3a\xb8\x17\x7f\x3d\x30\xff\x34\x2c\xe3\x70\x83\xf3\x39\xf9\xc3\xd6\x39\xe5\xe0\x34\xae\xe0\x6a\x1d\x0c\x31\xa4\xde\x2f\xe9\xd4\xd6\x7d\x38\x4a\x9f\x47\xe5\x7a\x98\xc9\x42\x1a\xfd\x38\x48\x8a\x1c\xe8\xca\xea\x72\xd9\xb5\x6d\xc8\x8c\x3e\x6e\xff\xe7\x38\x50\xb2\x3d\xd3\x92\xb6\x7c\xd2\x77\x6a\xfa\x66\x6e\xd5\xcf\x5e\x40\xad\x0a\xf7\x5f\x29\x3c\x9c\xbc\x44\xe0\x41\x99\xdb\x81\xec\x41\xd8\x41\x73\x71\xde\xf3\xc7\xd0\xc5\xf9\x60\xab\xa7\x28\xd9\x2c\x11\x18\x3e\x1d\x4d\x86\x83\x3f\x03\x00\x00\xff\xff\x50\x1f\xdf\x65\x98\x13\x00\x00")

func swaggerSwaggerYamlBytes() ([]byte, error) {
	return bindataRead(
		_swaggerSwaggerYaml,
		"swagger/swagger.yaml",
	)
}

func swaggerSwaggerYaml() (*asset, error) {
	bytes, err := swaggerSwaggerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger/swagger.yaml", size: 5016, mode: os.FileMode(420), modTime: time.Unix(1507300625, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"swagger/swagger.json": swaggerSwaggerJson,
	"swagger/swagger.yaml": swaggerSwaggerYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"swagger": &bintree{nil, map[string]*bintree{
		"swagger.json": &bintree{swaggerSwaggerJson, map[string]*bintree{}},
		"swagger.yaml": &bintree{swaggerSwaggerYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

