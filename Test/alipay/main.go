package main
import(
	"fmt"
	"github.com/smartwalle/alipay/v3"

)
func main()  {
	var appId = "2021000116679223"
	var privateKey = "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCBXrzfCsHzV2yBxn/KYO3i59djrM4ckMW6JLJyxkyAl8kzBs7MaMXa22NKlwrKC0AsUsI4g1BQvtbelBLCULrnEeosIlohOfIQrJwQpEu+2rCtvM2YVXS0h0h8MTE4kLO3wwMTCeusogKv0S5gNmyMgWb2eEJlv257sNA4iKyHRtFhpxk6X0vqpegf8ZZpaGKrPuC3JjpzJJfbI2oiU4/9YR4WeNBkfWyYSIIEeSPh/sDT5sOpOAoVGi1CzVgeMpSvhBWseNVhnVg34BmjbIe8XLleyw9tN3BLZMp/qfEz3/JcOmJOfHxoDss75aSUwtMkd5D8UMOe8cVHDYT90EwPAgMBAAECggEATXUysnyXaaJTdlQqGTr106CqZqSFd7b6nBvyuCOglHHM5n9R/DNTG1m11mge7p/T2XXnkbyVrvLEZdnUbYG2ljk0sx4SRsiR9YfTnWcxbuEzXaKomme4C4rgTHOLm+mPoRvi0FhlQiRyZWBiWvra/TlRM4sHfjIi7W4NDXxPiGJu7RYI4gv2B8OPCjXb0c8X+XPudrrXU0BlNsVt1Cb7DzlM0Wzv7+oVO6Z/QEBMeWiMOsbxIKNuPDgsczJ14eEyGWijwPBnQQvBYpkpfs1eLsL5jlDS+nRuBMz1ev7682mNrUYlenx4epxlNvMrB6z5dn2l6yRIzwf0q2yBz6cpAQKBgQDoYnpJ8iN1OYtjMiQnOR/RAlWVO2VRPjOB/gwMkSRLnS9n9+YXO77r9bImhhKNDkhDrk205vNgRbuwiehlIUWYiB0e3EffFs4rU5ZBjpikUwqnO9ZKdxZ0bONeC3t3ESnek8Gg8EdMBzo24yGhi0Se3aa2rZzvAEsRaeEEhR/AjwKBgQCOhFHO5Y6zdPKFacfw2eXs+XBsSfWgISfpPNDFRVEej4SPmlFXupjnK4EaOmPEGYwBkkTgpiv1EWxTtzmIrNrLvbGNp1e70fMqhxgmM5BPw188BP6wzn39GB0bX1IoaXFF978joVQJBjK3mCUfoHYMU46E9CchtOm7vhyAz0V8gQKBgQCnUJAfwZlKA56aUzW4j+aahAW+pr8yGYjYZWOjgLUTRB8nylL+E2RJW+Ni3VFqAgiBwnPsdgRxIoQafZC3j5ceVZIx8ARHWZIjm9EpblP0rF2VPv3xK4EdXnDt+3JvrgnpWZUmHRoYYXPGpQ5H05Aamhg4mxPM+PcTmJoMuRS9fQKBgQCOO/LluWdNJTdx/Rul/eIiOuR/vuScdtq9RYvahg6qoHOdWlc6Zil59Yo4ofO14AVCADgruMyAHm5tspyCEnmfA2fzxwKhmazFUeTBI48we/1NCnMiBEPrV6idC+oUGQAK43Jo3fbftsbhQsAyK1QFg1Lm0EFJu6PWUqthxvtDgQKBgQDmlfo6PKDMXyKI8WNSZdcAsQEp4MP0fVfJMB/ioC2w17d4XtvkvOFvPtFgvYo91vhorMSa4j2XG1M+fm+c1GLrDeCPWnOddjAL/saeO6EMWmKkpO4Xrc182Cm5v06PfdRN+VbTgo1x2AL/1Ai7xh96jHOuPbSqNjnqEg324n41EA==" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	//var aliPublicKey="MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkBui320sHC572qEBT8oCmaswX32x3YO2Di07XlkG71B+tdSJ7R6GzjVacrmGD2H9METl5b1Et0hbwrcBwnt/ER5Pyk1eprIy1UoTzcsJFH7sGfPNEET0jl0f3axwvFFNsIf9vi4mGjuR+YYZPcjYKpQ7ZhGT2I9mkHfKbkViQshe5LtlW9Qi4ag+vUl8vV0vkD5/OO1qP0c8zbwUmPpXO572TbuMVz4ZR+QSXsOUK/XTQbto3I2aBj12QnIwStOlH1WHxqnuxE3JQSYhS1p0jAb8dXTB2b2gdBwEDv9jahSNv1i1xRFPrnvH4tjkrn1frnJ2pEIICD1+dYg0624ztwIDAQAB"

	 client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		fmt.Println("实例化支付宝失败 失败")
		return
	}
	//client.LoadAliPayPublicKey(aliPublicKey)
	//
	err =client.LoadAppPublicCertFromFile("alipay/appCertPublicKey.crt") // 加载应用公钥证书
	// 将 key 的验证调整到初始化阶段
	if err != nil {
		fmt.Println("加载应用公钥证书 失败")
		return
	}
	err =client.LoadAliPayRootCertFromFile("alipay/alipayRootCert.crt") // 加载支付宝根证书
	if err != nil {
		fmt.Println("加载支付宝根证书 失败")
		return
	}
	err =client.LoadAliPayPublicCertFromFile("alipay/alipayCertPublicKey_RSA2.crt") // 加载支付宝公钥证书
	if err != nil {
		fmt.Println("加载支付宝公钥证书 失败")
		return
	}
	// 将 key 的验证调整到初始化阶段
	if err != nil {
		fmt.Println(err)
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://www.code688.com"
	p.ReturnURL = "http://www.code688.com"
	p.Subject = "测试"
	p.OutTradeNo = "333333333333111"
	p.TotalAmount = "11.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	 url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()
	fmt.Println(payURL)
}

