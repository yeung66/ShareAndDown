package utils

import (
	"fmt"
	qrcode "github.com/yeqown/go-qrcode"
)

func GenQRCode(fileUrl string, codePath string) error {
	qrc, err := qrcode.New(fileUrl)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return err
	}

	// save file
	if err = qrc.Save(codePath); err != nil {
		fmt.Printf("could not save image: %v", err)
	}

	return err
}
