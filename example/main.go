
package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/Ma-Chenwei/cryptex"
)

func main() {
    fmt.Println("================================")
    fmt.Println("       Cryptex v1.1 Example")
    fmt.Println("================================")

    password := "MyMasterPassword"

    // =================================
    // 1. 自动生成 C 并加密
    // =================================

    fmt.Println("\n[1] Automatic C generation")

    original := []byte("Hello Cryptex!")

    c, saltA, nonceA, sCipher,
        saltB, nonceB, vaultCipher, err :=
        cryptex.Encrypt(password, nil, original)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("C length:", len(c))
    fmt.Println("Salt A:", len(saltA))
    fmt.Println("Nonce A:", len(nonceA))
    fmt.Println("System Cipher:", len(sCipher))
    fmt.Println("Salt B:", len(saltB))
    fmt.Println("Nonce B:", len(nonceB))
    fmt.Println("Vault Cipher:", len(vaultCipher))

    // =================================
    // 2. 解密
    // =================================

    fmt.Println("\n[2] Decryption")

    decrypted, err := cryptex.Decrypt(
        password,
        saltA,
        nonceA,
        sCipher,
        saltB,
        nonceB,
        vaultCipher,
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Original:", string(original))
    fmt.Println("Decrypted:", string(decrypted))
    fmt.Println("Match:", bytes.Equal(original, decrypted))

    // =================================
    // 3. 手动生成 C
    // =================================

    fmt.Println("\n[3] Custom C generation")

    customC, err := cryptex.GenerateC()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Generated C length:", len(customC))

    actualC, sa, na, sc, sb, nb, vc, err :=
        cryptex.Encrypt(password, customC, original)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Custom C match:", bytes.Equal(customC, actualC))

    // =================================
    // 4. 使用手动 C 的密文解密
    // =================================

    fmt.Println("\n[4] Custom C decryption")

    result, err := cryptex.Decrypt(
        password,
        sa, na, sc,
        sb, nb, vc,
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Decrypted:", string(result))

    // =================================
    // 5. 二进制数据加密
    // =================================

    fmt.Println("\n[5] Binary data encryption")

    binaryData := []byte{
        0x00, 0x01, 0x02, 0x03,
        0xFF, 0xFE, 0xFD, 0xFC,
    }

    _, bsa, bna, bsc, bsb, bnb, bvc, err :=
        cryptex.Encrypt(password, nil, binaryData)

    if err != nil {
        log.Fatal(err)
    }

    binaryResult, err := cryptex.Decrypt(
        password,
        bsa, bna, bsc,
        bsb, bnb, bvc,
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Binary original:", binaryData)
    fmt.Println("Binary decrypted:", binaryResult)
    fmt.Println("Binary match:", bytes.Equal(binaryData, binaryResult))

    fmt.Println("\n================================")
    fmt.Println("       Example completed")
    fmt.Println("================================")
}
