package main

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func xorBytes(a, b []byte) []byte {

	out := make([]byte, len(a))

	for i := range a {
		out[i] = a[i] ^ b[i]
	}

	return out
}

func calculateOPc(kiHex string, opHex string) (string, error) {

	kiHex = strings.ReplaceAll(kiHex, " ", "")
	opHex = strings.ReplaceAll(opHex, " ", "")

	ki, err := hex.DecodeString(kiHex)

	if err != nil {
		return "", fmt.Errorf("Invalid Ki")
	}

	op, err := hex.DecodeString(opHex)

	if err != nil {
		return "", fmt.Errorf("Invalid OP")
	}

	if len(ki) != 16 {
		return "", fmt.Errorf("Ki must be 16 bytes")
	}

	if len(op) != 16 {
		return "", fmt.Errorf("OP must be 16 bytes")
	}

	block, err := aes.NewCipher(ki)

	if err != nil {
		return "", err
	}

	encrypted := make([]byte, 16)

	block.Encrypt(encrypted, op)

	opc := xorBytes(encrypted, op)

	return strings.ToUpper(hex.EncodeToString(opc)), nil
}

func main() {

	a := app.New()

	w := a.NewWindow("SIM OPc Calculator")

	kiEntry := widget.NewEntry()
	kiEntry.SetPlaceHolder("Ki (32 hex chars)")

	opEntry := widget.NewEntry()
	opEntry.SetPlaceHolder("OP (32 hex chars)")

	result := widget.NewEntry()
	result.Disable()

	status := widget.NewLabel("")

	calc := widget.NewButton("Calculate OPc", func() {

		opc, err := calculateOPc(kiEntry.Text, opEntry.Text)

		if err != nil {

			status.SetText(err.Error())

			return
		}

		result.SetText(opc)

		status.SetText("OK")
	})

	copy := widget.NewButton("Copy", func() {

		w.Clipboard().SetContent(result.Text)

		status.SetText("Copied")
	})

	form := container.NewVBox(

		widget.NewLabel("Ki"),
		kiEntry,

		widget.NewLabel("OP"),
		opEntry,

		calc,

		widget.NewLabel("OPc"),
		result,

		copy,

		status,
	)

	w.SetContent(form)

	w.Resize(fyne.NewSize(420, 300))

	w.ShowAndRun()
}
