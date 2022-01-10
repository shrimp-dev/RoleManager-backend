package pix

import "fmt"

const (
	FORMAT             = "000201"
	GUI                = "0014br.gov.bcb.pix"
	CATEGORY           = "0000"
	ADITIONAL_TEMPLATE = "0503***"
)

const (
	ACCOUNT_OPCODE   = 26
	CATEGORY_OPCODE  = 52
	CURRENCY_OPCODE  = 53
	AMOUNT_OPCODE    = 54
	COUNTRY_OPCODE   = 58
	NAME_OPCODE      = 59
	CITY_OPCODE      = 60
	ADITIONAL_OPCODE = 62
	CRC_OPCODE       = 63
)

const (
	BRL = "986"
)

const (
	BRAZIL = "BR"
)

type Pix struct {
	PixKey      string
	ReciverName string
	ReciverCity string
	Amount      float32
	Info        string
}

func (p *Pix) GeneratePixStream() string {
	message := FORMAT
	message += p.account()
	message += p.category()
	message += p.currency()
	message += p.amount()
	message += p.country()
	message += p.name()
	message += p.city()
	message += p.aditional()
	message += p.crc16chk(message)

	return message
}

func (p *Pix) account() string {
	accPayload := GUI + p.PixKey + p.Info

	accStr := getOpCodeStr(ACCOUNT_OPCODE)
	accStr += parseStringToEmv(accPayload)

	return accStr
}

func (p *Pix) category() string {
	catStr := getOpCodeStr(CATEGORY_OPCODE)
	catStr += parseStringToEmv(CATEGORY)

	return catStr
}

func (p *Pix) currency() string {
	crStr := getOpCodeStr(CURRENCY_OPCODE)
	crStr += parseStringToEmv(BRL)

	return crStr
}

func (p *Pix) amount() string {
	amountStr := fmt.Sprint(p.Amount)

	amtStr := getOpCodeStr(AMOUNT_OPCODE)
	amtStr += parseStringToEmv(amountStr)

	return amtStr
}

func (p *Pix) country() string {
	countryStr := getOpCodeStr(COUNTRY_OPCODE)
	countryStr += parseStringToEmv(BRAZIL)

	return countryStr
}

func (p *Pix) name() string {
	name := p.ReciverName
	if len(name) > 25 {
		name = name[:25]
	}

	nameStr := getOpCodeStr(NAME_OPCODE)
	nameStr += parseStringToEmv(name)

	return nameStr
}

func (p *Pix) city() string {
	city := p.ReciverCity
	if len(city) > 25 {
		city = city[:25]
	}

	cityStr := getOpCodeStr(CITY_OPCODE)
	cityStr += parseStringToEmv(city)

	return cityStr
}

func (p *Pix) aditional() string {
	aditionalStr := getOpCodeStr(ADITIONAL_OPCODE)
	aditionalStr += parseStringToEmv(ADITIONAL_TEMPLATE)

	return aditionalStr
}

func (p *Pix) crc16chk(msg string) string {
	msg += getOpCodeStr(CRC_OPCODE)
	msg += contentLen("ffff")

	chk := Calculate_CRC_CCITT(msg)
	chkStr := fmt.Sprintf("%04x", chk)

	crcStr := getOpCodeStr(CRC_OPCODE)
	crcStr += parseStringToEmv(chkStr)

	return crcStr
}

func parseStringToEmv(str string) string {
	res := contentLen(str)
	res += str

	return res
}

func contentLen(i string) string {
	len := len(i) - 1
	str := fmt.Sprintf("%02d", len)
	return str
}

func getOpCodeStr(opcode int) string {
	return fmt.Sprint(opcode)
}
