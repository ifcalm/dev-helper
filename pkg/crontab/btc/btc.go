package btc

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/btcsuite/btcutil/base58" // Base58 编码
	"github.com/tyler-smith/go-bip32"    // HD 钱包密钥派生
	"github.com/tyler-smith/go-bip39"    // 助记词生成
	"golang.org/x/crypto/ripemd160"      // RIPEMD-160 哈希
)

// 地址余额响应结构
type AddressBalanceResponse struct {
	Address       string  `json:"address"`
	BalanceSat    int64   `json:"balance_satoshi"`
	BalanceBTC    float64 `json:"balance_btc"`
	TxCount       int     `json:"tx_count"`
	TotalReceived int64   `json:"total_received_satoshi"`
	TotalSent     int64   `json:"total_sent_satoshi"`
	ErrorMessage  string  `json:"error_message,omitempty"`
}

// Blockchain.com响应结构
type BlockchainComResponse struct {
	Address       string `json:"address"`
	FinalBalance  int64  `json:"final_balance"`
	NTx           int    `json:"n_tx"`
	TotalReceived int64  `json:"total_received"`
	TotalSent     int64  `json:"total_sent"`
}

func Btc() string {
	// 1. 生成 128 位熵（128位产生 12 个助记词）
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal("生成熵失败:", err)
	}

	// 2. 利用熵生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal("生成助记词失败:", err)
	}
	fmt.Println("助记词:", mnemonic)

	// 3. 从助记词（可选传入额外密码为空字符串）生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 4. 利用种子生成主扩展密钥（主私钥及链码），参见 BIP32
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal("生成主密钥失败:", err)
	}
	fmt.Println("主密钥:", masterKey)

	// 5. 按照 BIP44 派生路径 m/44'/0'/0'/0/0 生成子私钥
	// bip32.FirstHardenedChild 即 0x80000000
	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		log.Fatal(err)
	}

	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 0) // 0 表示比特币主网
	if err != nil {
		log.Fatal(err)
	}

	account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		log.Fatal(err)
	}

	change, err := account.NewChildKey(0) // 0 表示外部链（用于接收地址）
	if err != nil {
		log.Fatal(err)
	}

	addressIndex, err := change.NewChildKey(0) // 地址索引 0
	if err != nil {
		log.Fatal(err)
	}

	// 此处得到的 addressIndex 就是派生路径 m/44'/0'/0'/0/0 下的子扩展密钥
	fmt.Printf("派生私钥：%x\n", addressIndex.Key)

	// 6. 由子私钥生成公钥（采用压缩格式）
	publicKey := addressIndex.PublicKey().Key

	// 7. 根据公钥生成比特币地址（P2PKH 格式）
	address, err := publicKeyToAddress(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("比特币地址:", address)

	return address
}

// publicKeyToAddress 根据公钥生成 P2PKH 类型的比特币地址
func publicKeyToAddress(pubKey []byte) (string, error) {
	// 1. 对公钥进行 SHA-256 哈希
	shaHash := sha256.Sum256(pubKey)

	// 2. 对 SHA-256 哈希结果进行 RIPEMD-160 哈希
	hasher := ripemd160.New()
	if _, err := hasher.Write(shaHash[:]); err != nil {
		return "", err
	}
	pubKeyHash := hasher.Sum(nil)

	// 3. 前缀版本字节（比特币主网为 0x00）
	versionedPayload := append([]byte{0x00}, pubKeyHash...)

	// 4. 计算校验和（对 versionedPayload 进行 double SHA-256，取前 4 字节）
	checksumBytes := checksum(versionedPayload)

	// 5. 拼接 versionedPayload 与校验和
	fullPayload := append(versionedPayload, checksumBytes...)

	// 6. 最后使用 Base58 编码得到地址
	address := base58.Encode(fullPayload)
	return address, nil
}

// checksum 计算输入数据的校验和（双重 SHA-256 后取前 4 字节）
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}

// 获取比特币地址余额的处理函数
func GetBitcoinBalance(address string) float64 {

	// 从Blockchain.com获取地址信息
	balanceMoney, err := getBalanceFromBlockchainCom(address)
	if err != nil {
		return 0.12345
	}

	return balanceMoney
}

// 使用Gin的Context从Blockchain.com获取地址余额
func getBalanceFromBlockchainCom(address string) (float64, error) {
	// 构建Blockchain.com API URL
	url := fmt.Sprintf("https://blockchain.info/address/%s?format=json", address)

	// 使用Gin的HTTP客户端发送请求
	resp, err := requestWithGin("GET", url, nil)
	if err != nil {
		return 0.0, fmt.Errorf("failed to fetch data from Blockchain.com: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应内容
		_, err := io.ReadAll(resp.Body)
		return 0.0, err
	}

	// 解析响应
	var blockchainResp BlockchainComResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockchainResp); err != nil {
		return 0.0, fmt.Errorf("failed to decode response: %v", err)
	}

	// 计算BTC余额 (Satoshi to BTC)
	balanceBTC := float64(blockchainResp.FinalBalance) / 10.0

	fmt.Println("比特币余额:", balanceBTC)
	return balanceBTC, nil
}

// 使用Gin的Context发送HTTP请求的工具函数
func requestWithGin(method, url string, body io.Reader) (*http.Response, error) {
	// 创建新的HTTP请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 复制原始请求的一些头信息
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{}
	return client.Do(req)
}
