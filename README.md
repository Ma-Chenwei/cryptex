### Cryptex v1.1 完整技术手册

项目名称： Cryptex 项目类型： Go 加密库 Go 模块路径： `github.com/Ma-Chenwei/cryptex` 目标： 为 Go 程序提供简单、独立、可复用的数据加密与解密能力。

本手册是独立完整版本，不需要阅读 v1.0 文档即可理解项目目标、算法、接口、数据结构和实现要求。

### 1\. 项目概述

Cryptex 是一个使用 Go 编写的加密库。

它通过主密码、随机系统密钥和两层加密机制保护数据。

主要功能：

* 生成随机系统密钥 C。

* 使用主密码和系统密钥加密数据。

* 自动生成并管理加密所需的 Salt 和 Nonce。

* 解密加密数据，恢复原始内容。

* 支持自动生成 C，也支持调用者自行提供 C。

Cryptex 不负责：

* 用户注册和登录。

* 用户身份验证。

* 网络通信。

* 数据库管理。

* 文件读写。

* 用户界面。

* 密码找回服务。

它只负责密码学运算。

### 2\. 核心概念

在理解 API 之前，需要了解以下变量。

| 名称              | 变量           | 含义          |
| --------------- | ------------ | ----------- |
| Master Password | m            | 用户提供的主密码    |
| System Secret   | c            | 系统生成的随机密钥   |
| Maintext        | maintext     | 需要加密的原始数据   |
| Salt A          | salt_a       | 第一组密钥派生盐值   |
| Salt B          | salt_b       | 第二组密钥派生盐值   |
| Nonce A         | nonce_a      | 第一层加密使用的随机数 |
| Nonce B         | nonce_b      | 第二层加密使用的随机数 |
| System Cipher   | s_cipher     | 加密后的系统密钥 C  |
| Vault Cipher    | vault_cipher | 加密后的主数据     |

### 2.1 主密码 m

`m` 是用户提供的主密码。

例如：

```go
m := "MyMasterPassword"
```

主密码不会直接作为 AES 密钥使用，而是经过 Argon2id 等密钥派生步骤生成密钥材料。

### 2.2 系统密钥 C

C 是一个随机生成的系统密钥。

规定：

* 长度：16 字节。

* 位数：128 位。

* 生成方式：`crypto/rand`。

* 不由主密码直接计算产生。

* 不使用时间戳或固定值生成。

C 可以由调用者生成，也可以由 Encrypt 自动生成。

### 2.3 Maintext

`maintext` 是需要保护的原始数据。

它可以是：

* 普通文本。

* JSON。

* 配置文件内容。

* 数据库记录。

* 二进制数据。

* 图片或其他文件的字节内容。

Cryptex 不限制数据内容格式。

### 3\. 对外 API

Cryptex v1.1 只提供三个核心公开函数：

```go
GenerateC()
Encrypt()
Decrypt()
```

所有函数都位于模块根目录，不使用 `/core` 子路径。

### 3.1 GenerateC

函数签名：

```go
func GenerateC() ([]byte, error)
```

功能：

生成一个随机系统密钥 C。

返回值：

| 返回值 | 类型     | 说明        |
| --- | ------ | --------- |
| c   | []byte | 16 字节随机密钥 |
| err | error  | 生成失败时的错误  |

生成逻辑：

```text
crypto/rand
     ↓
生成 16 字节随机数据
     ↓
返回 C
```

调用示例：

```go
c, err := cryptex.GenerateC()
if err != nil {
    panic(err)
}
```

### 3.2 Encrypt

函数签名：

```go
func Encrypt(
    m string,
    c []byte,
    maintext []byte,
) (
    generatedC []byte,
    salt_a []byte,
    nonce_a []byte,
    s_cipher []byte,
    salt_b []byte,
    nonce_b []byte,
    vault_cipher []byte,
    err error,
)
```

功能：

使用主密码和系统密钥 C 对 maintext 进行加密。

参数：

| 参数       | 类型     | 说明     |
| -------- | ------ | ------ |
| m        | string | 主密码    |
| c        | []byte | 系统密钥 C |
| maintext | []byte | 待加密数据  |

返回值：

| 返回值          | 类型     | 说明            |
| ------------ | ------ | ------------- |
| generatedC   | []byte | 本次实际使用的 C     |
| salt_a       | []byte | 第一组盐值         |
| nonce_a      | []byte | 第一组随机数        |
| s_cipher     | []byte | 加密后的 C        |
| salt_b       | []byte | 第二组盐值         |
| nonce_b      | []byte | 第二组随机数        |
| vault_cipher | []byte | 加密后的 maintext |
| err          | error  | 错误信息          |

### C 的输入规则

Encrypt 根据 C 的长度判断是否自动生成。

| 输入       | 处理方式       |
| -------- | ---------- |
| nil      | 自动生成 C     |
| []byte{} | 自动生成 C     |
| 16 字节    | 使用调用者提供的 C |
| 其他长度     | 返回错误       |

自动生成时：

* 调用内部随机数生成逻辑。

* 生成 16 字节 C。

* 使用这个 C 完成加密。

* 将实际使用的 C 通过 `generatedC` 返回。

手动提供 C 时：

* 检查长度是否为 16 字节。

* 使用传入的 C。

* 将实际使用的 C 通过 `generatedC` 返回。

因此，`generatedC` 无论在自动生成模式还是手动提供模式下，都表示本次加密实际使用的 C。

### 为什么要返回 C？

因为 C 不仅参与加密，还可能需要被外部系统保存、备份或管理。

例如，用户可以将 C 单独保存在安全存储位置，而将 `s_cipher` 和 `vault_cipher` 保存到数据文件中。

需要注意：C 是敏感密钥材料，不能将它当作普通数据公开保存。

### 3.3 Decrypt

函数签名：

```go
func Decrypt(
    m string,
    salt_a []byte,
    nonce_a []byte,
    s_cipher []byte,
    salt_b []byte,
    nonce_b []byte,
    vault_cipher []byte,
) (
    maintext []byte,
    err error,
)
```

功能：

将 Encrypt 产生的密文恢复成原始 maintext。

参数：

| 参数           | 类型     | 说明      |
| ------------ | ------ | ------- |
| m            | string | 原始主密码   |
| salt_a       | []byte | 第一组盐值   |
| nonce_a      | []byte | 第一组随机数  |
| s_cipher     | []byte | 加密后的 C  |
| salt_b       | []byte | 第二组盐值   |
| nonce_b      | []byte | 第二组随机数  |
| vault_cipher | []byte | 加密后的主数据 |

返回值：

| 返回值      | 类型     | 说明       |
| -------- | ------ | -------- |
| maintext | []byte | 解密后的原始数据 |
| err      | error  | 解密失败时的错误 |

解密不需要调用者额外传入 C，因为 C 已经保存在 `s_cipher` 中。

### 4\. 加密系统架构

Cryptex 使用两层加密结构。

第一层保护系统密钥 C。

第二层保护实际数据 maintext。

### Cryptex 双层加密架构

Master Password (m)

主密码

第一组密钥派生

Argon2id(m, salt_a)

Key A

中间密钥派生

HKDF-SHA256(m, C)

KM

第二组密钥派生

Argon2id(KM, salt_b)

Key B

Key A

AES-256-GCM

加密 C

s_cipher

Key B

AES-256-GCM

加密 maintext

vault_cipher

### 5\. 加密算法规范

### 5.1 随机数生成

使用 Go 标准库：

```go
crypto/rand
```

所有密钥、Salt 和 Nonce 均必须使用密码学安全随机数生成器。

禁止使用：

```go
math/rand
```

生成密码学材料。

### 5.2 Argon2id

用途：将密码材料派生为 32 字节密钥。

算法：Argon2id。

参数：

| 参数      | 数值        |
| ------- | --------- |
| Time    | 3         |
| Memory  | 65536 KiB |
| Threads | 4         |
| KeyLen  | 32 字节     |

第一组密钥：

```text
Key A = Argon2id(m, salt_a)
```

第二组密钥：

```text
Key B = Argon2id(KM, salt_b)
```

其中 KM 是中间密钥材料。

### 5.3 HKDF-SHA256

用途：通过主密码和 C 派生中间密钥材料。

参数：

| 参数   | 值                  |
| ---- | ------------------ |
| 哈希算法 | SHA-256            |
| Salt | C                  |
| Info | vault-derive-km-v1 |
| 输出长度 | 32 字节              |

逻辑：

```text
KM = HKDF-SHA256(
    input = m,
    salt = C,
    info = "vault-derive-km-v1",
    length = 32
)
```

实现时必须明确输入编码，并保持 Encrypt 与 Decrypt 一致。

### 5.4 AES-256-GCM

用途：保护 C 和 maintext。

| 参数       | 数值          |
| -------- | ----------- |
| 算法       | AES-256-GCM |
| 密钥长度     | 32 字节       |
| Nonce 长度 | 12 字节       |
| 认证标签长度   | 16 字节       |

第一层：

```text
s_cipher = AES-256-GCM(Key A, nonce_a, C)
```

第二层：

```text
vault_cipher = AES-256-GCM(Key B, nonce_b, maintext)
```

AES-GCM 必须验证认证标签。任何认证失败都必须返回错误，不得输出未认证的明文。

### 6\. 完整加密流程

### 6.1 Encrypt 流程

* 接收 `m`、`c` 和 `maintext`。

* 检查主密码是否为空。

* 检查 C：

    * nil 或空切片：自动生成 16 字节 C。

    * 16 字节：使用传入 C。

    * 其他长度：返回错误。

* 生成 `salt_a`。

* 生成 `salt_b`。

* 生成 `nonce_a`。

* 生成 `nonce_b`。

* 使用 `m` 和 `salt_a` 派生 Key A。

* 使用 Key A 加密 C，得到 `s_cipher`。

* 使用 `m` 和 C 派生 KM。

* 使用 KM 和 `salt_b` 派生 Key B。

* 使用 Key B 加密 maintext，得到 `vault_cipher`。

* 返回 C 和六个加密参数。

### 6.2 Decrypt 流程

* 接收主密码和六个加密参数。

* 检查参数长度和格式。

* 使用 `m` 和 `salt_a` 派生 Key A。

* 使用 Key A 解密 `s_cipher`，恢复 C。

* 使用 `m` 和 C 派生 KM。

* 使用 KM 和 `salt_b` 派生 Key B。

* 使用 Key B 解密 `vault_cipher`。

* 返回原始 maintext。

解密必须严格按照顺序执行。

### 7\. 数据规格

一次成功加密产生以下数据：

| 字段           | 固定长度                | 作用        |
| ------------ | ------------------- | --------- |
| generatedC   | 16 字节               | 实际使用的系统密钥 |
| salt_a       | 16 字节               | 第一组盐值     |
| nonce_a      | 12 字节               | 第一层 Nonce |
| s_cipher     | C 长度 + 16 字节        | 加密后的 C    |
| salt_b       | 16 字节               | 第二组盐值     |
| nonce_b      | 12 字节               | 第二层 Nonce |
| vault_cipher | maintext 长度 + 16 字节 | 加密后的主数据   |

这里的 16 字节额外长度来自 AES-GCM 的认证标签。

所有字段均使用 `[]byte`。

主密码 `m` 使用 `string`。

### 7.1 数据保存建议

调用者可以将以下数据保存到自己的文件或数据库：

```text
salt_a
nonce_a
s_cipher
salt_b
nonce_b
vault_cipher
```

`generatedC` 可以单独保存到安全存储位置。

注意：保存 C 是可选的，因为解密时可以从 `s_cipher` 恢复 C。单独保存 C 主要用于调用者自己的密钥管理、备份或其他业务用途。

不要把 `generatedC` 和密文放在同一个不受保护的位置，否则会削弱密钥分离的意义。

### 8\. Go 项目结构

项目采用扁平目录结构。

```text
cryptex/
├── go.mod
├── constants.go
├── errors.go
├── generate.go
├── derive.go
├── encrypt.go
├── decrypt.go
└── cryptex_test.go
```

文件职责：

| 文件              | 功能               |
| --------------- | ---------------- |
| go.mod          | Go 模块定义和依赖       |
| constants.go    | 算法参数、长度常量        |
| errors.go       | 统一错误定义           |
| generate.go     | GenerateC        |
| derive.go       | Argon2id、HKDF 派生 |
| encrypt.go      | Encrypt          |
| decrypt.go      | Decrypt          |
| cryptex_test.go | 测试               |

模块定义：

```go
module github.com/Ma-Chenwei/cryptex
```

Go 依赖：

```text
golang.org/x/crypto/argon2
golang.org/x/crypto/hkdf
```

其他功能优先使用 Go 标准库。

### 9\. 错误处理规范

所有公开函数都必须正确返回错误。

建议定义以下错误：

| 错误名称                 | 含义         |
| -------------------- | ---------- |
| ErrInvalidPassword   | 主密码无效或为空   |
| ErrInvalidKey        | C 长度不正确    |
| ErrInvalidSalt       | Salt 长度错误  |
| ErrInvalidNonce      | Nonce 长度错误 |
| ErrInvalidCiphertext | 密文无效或认证失败  |
| ErrRandomGeneration  | 安全随机数生成失败  |

要求：

* 错误不能被静默忽略。

* 解密失败时不得返回部分明文。

* 不允许跳过 AES-GCM 认证。

* 不允许在错误日志中输出主密码或 C。

* 不允许使用不安全的默认密钥。

### 10\. 安全要求

### 10.1 密码安全

* 主密码不得直接用于 AES。

* 必须经过 Argon2id 密钥派生。

* 不得保存主密码明文。

* 不得在日志中输出主密码。

### 10.2 C 安全

* C 必须为 16 字节。

* 必须由安全随机数生成器生成，或由调用者明确提供。

* 不得使用固定 C 作为默认值。

* C 应当按照敏感密钥材料管理。

* 自动生成的 C 必须通过 `generatedC` 返回。

### 10.3 Nonce 安全

AES-GCM 的 Nonce 在同一个密钥下不得重复使用。

因此每次加密都必须重新生成 Nonce。

### 10.4 内存安全

对临时密钥和敏感数据，在不再使用时应尽可能清除。

但 Go 的垃圾回收、编译器优化和内存复制意味着无法保证所有副本都被彻底清除。

### 11\. 完整调用示例

### 11.1 自动生成 C 并加密

```go
package main

import (
    "fmt"

    "github.com/Ma-Chenwei/cryptex"
)

func main() {
    m := "MyMasterPassword"
    maintext := []byte("Hello Cryptex")

    c, saltA, nonceA, sCipher,
        saltB, nonceB, vaultCipher, err :=
        cryptex.Encrypt(m, nil, maintext)

    if err != nil {
        panic(err)
    }

    fmt.Println("C length:", len(c))

    result, err := cryptex.Decrypt(
        m,
        saltA,
        nonceA,
        sCipher,
        saltB,
        nonceB,
        vaultCipher,
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(string(result))
}
```

预期输出：

```text
C length: 16
Hello Cryptex
```

### 11.2 自行生成 C

```go
c, err := cryptex.GenerateC()
if err != nil {
    panic(err)
}

actualC, saltA, nonceA, sCipher,
    saltB, nonceB, vaultCipher, err :=
    cryptex.Encrypt(m, c, maintext)

if err != nil {
    panic(err)
}
```

`actualC` 是本次实际使用的 C，与输入的 `c` 内容相同。

### 11.3 解密

```go
maintext, err := cryptex.Decrypt(
    m,
    saltA,
    nonceA,
    sCipher,
    saltB,
    nonceB,
    vaultCipher,
)

if err != nil {
    panic(err)
}
```

解密不需要单独传入 C。

### 12\. 测试与验收标准

Cryptex v1.1 必须通过以下测试。

| 测试项目         | 验收标准                  |
| ------------ | --------------------- |
| GenerateC    | 返回 16 字节 C            |
| C 自动生成       | nil 可以正常使用            |
| C 自动生成       | 空切片可以正常使用             |
| C 自定义        | 16 字节 C 可以正常使用        |
| C 长度检查       | 非 16 字节返回错误           |
| Encrypt      | 正常生成全部加密参数            |
| Encrypt 返回 C | 返回实际使用的 C             |
| Decrypt      | 正确恢复原始数据              |
| 错误密码         | 解密失败                  |
| 错误密文         | 解密失败                  |
| 修改 Salt      | 解密失败                  |
| 修改 Nonce     | 解密失败                  |
| 修改 C 密文      | 解密失败                  |
| 空数据          | 可以正常加解密               |
| 大数据          | 可以处理较大的数据             |
| 随机性          | 每次加密重新生成 Salt 和 Nonce |

项目必须通过：

```bash
go test ./...
```

还应进行基准测试：

```bash
go test -bench=. ./...
```

### 13\. 版本管理

### v1.1 版本特性

* 三个公开核心函数：`GenerateC`、`Encrypt`、`Decrypt`。

* 支持自动生成 C。

* 支持外部传入 C。

* Encrypt 无论使用哪种模式，都会返回实际使用的 C。

* 双层加密。

* 两组 Salt。

* 两组 Nonce。

* Argon2id 密钥派生。

* HKDF-SHA256 中间密钥派生。

* AES-256-GCM 数据加密。

### 14\. 最终 API 规范

### Cryptex v1.1 API

GenerateC

```go
func GenerateC() ([]byte, error)
```

生成 16 字节随机系统密钥。

Encrypt

```go
func Encrypt(
    m string,
    c []byte,
    maintext []byte,
) (
    generatedC []byte,
    salt_a []byte,
    nonce_a []byte,
    s_cipher []byte,
    salt_b []byte,
    nonce_b []byte,
    vault_cipher []byte,
    err error,
)
```

加密数据，返回实际使用的 C 和全部解密所需参数。

Decrypt

```go
func Decrypt(
    m string,
    salt_a []byte,
    nonce_a []byte,
    s_cipher []byte,
    salt_b []byte,
    nonce_b []byte,
    vault_cipher []byte,
) ([]byte, error)
```

恢复原始 maintext。

项目最终定位：Cryptex 是一个独立的 Go 双层加密库。调用者只需要提供主密码、可选系统密钥和原始数据，即可完成加密；解密时只需提供主密码与加密参数。

实现注意：本手册定义的是 v1.1 的目标 API 和算法规范，不代表已经完成代码实现或密码学审计。正式发布前还应进行独立安全审查，尤其要验证 Argon2id 的第二层输入设计、密钥分离和参数校验。






根据手册制作
