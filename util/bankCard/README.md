


## API List

- 获取所有支行 BankBranchList()
- 根据 areaID 获取当前区域下所有支行 BankListByArea(bankID string, areaID string)
- 检测是否是银行卡 IsBankCard(bankCardNo string)
- 根据卡号获取银行信息 GetBankByCardBin(bankCardNo string, bank *Bank)
- 使用阿里接口查询银行卡信息 GetBankByCardOnline(cardNo string, bankInfo *Bank)


## Usage
```golang
func main() {
    cardNo := "XXXXXX"
	var bank bankCard.Bank
	if err := bankCard.GetBankByCardOnline(cardNo, &bank); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("%+v\n", bank)
	if err := bankCard.GetBankByCardBin(cardNo, &bank); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("%+v\n", bank)
}

```

