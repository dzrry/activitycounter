package utilities

import (
	"bytes"
	"encoding/base64"

	"github.com/dzrry/activitycounter/insta"
)

// ExportAsBytes exports selected *Instagram object as []byte
func ExportAsBytes(inst *insta.Instagram) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := insta.Export(inst, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ExportAsBase64String exports selected *Instagram object as base64 encoded string
func ExportAsBase64String(inst *insta.Instagram) (string, error) {
	bytes, err := ExportAsBytes(inst)
	if err != nil {
		return "", err
	}

	sEnc := base64.StdEncoding.EncodeToString(bytes)
	return sEnc, nil
}
