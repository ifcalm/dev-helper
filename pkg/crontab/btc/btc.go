package btc

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58" // Base58 编码
	"github.com/tyler-smith/go-bip32"    // HD 钱包密钥派生
	"github.com/tyler-smith/go-bip39"    // 助记词生成
	"golang.org/x/crypto/ripemd160"      // RIPEMD-160 哈希
)

func btc() {
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
