package pix_test

import (
	"drinkBack/payments/pix"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrcChecksum(t *testing.T) {

	msg := "00020104141234567890123426580014BR.GOV.BCB.PIX0136123e4567-e12b-12d1-a456-42665544000027300012BR.COM.OUTRO011001234567895204000053039865406123.455802BR5917NOME DO RECEBEDOR6008BRASILIA61087007490062190515RP12345678-201980390012BR.COM.OUTRO01190123.ABCD.3456.WXYZ6304"
	chk := pix.Calculate_CRC_CCITT(msg)

	t.Logf("%4x", chk)
	assert.Equal(t, uint16(0xad38), chk)

}

func TestPix(t *testing.T) {
	px := pix.Pix{PixKey: "", ReciverName: "", ReciverCity: "FORTALEZA", Amount: 2.34}
	str := px.GeneratePixStream()
	t.Log(str)
}
