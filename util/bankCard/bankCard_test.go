package bankCard

import (
	"fmt"
	"log"
	"testing"
)

func TestBank_GetInfo(t *testing.T) {
	cardNo := "6226622219970718"
	var bank Bank
	if err := GetBankByCardOnline(cardNo, &bank); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("%+v\n", bank)
	if err := GetBankByCardBin(cardNo, &bank); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("%+v\n", bank)
}

func TestBankListByArea(t *testing.T) {
	area := BankListByArea("CEB", "1101")
	fmt.Println(area)
}

func TestBankBranchList(t *testing.T) {
	list := BankBranchList()
	fmt.Println(list)
}

func TestGetBankInfo(t *testing.T) {
	res, err := GetBankInfo("6226622219970718")
	fmt.Printf("%+v,%+v", res, err)
}
