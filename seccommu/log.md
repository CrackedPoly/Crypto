# openssl使用步骤

## 生成证书

### 生成密钥
`openssl genrsa -des3 -out ca.key 1024`

### 设置密钥密码为空
`openssl rsa -in ca.key -out ca.key`

### 生成根证书
`openssl req -new -x509 -key ca.key -out ca.crt -days 365`

### 生成子证书请求文件
`openssl req -new -key child.key -out child.csr`

### 生成子证书
`openssl x509 -req -in child.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out child.crt -days 365`

## 格式转换

### 证书格式转换
`openssl x509 -inform DER -in certificate.cer -out certificate.crt`

# 本次实验Windows下证书的使用

## 生成证书

### 生成根证书
`makecert -n "CN=Root" -r -sv RootIssuer.pvk RootIssuer.cer`

然后在图形化窗口中输入私钥的保护口令

### 使用根证书签发子证书
`makecert -n "CN=Child" -iv RootIssuer.pvk -ic RootIssuer.cer -sv ChildSubject.pvk ChildSubject.cer`

## 格式转换

### 将私钥和证书合为PKCS12格式的pfx
`pvk2pfx -pvk ChildSubject.pvk -spc ChildSubject.cer -pfx ChildSubject.pfx -pi test -po 666666 -f`

-pi指定根证书的私钥保护口令，-po指定pfx证书文件的保护口令(非常重要)，-f表示覆盖文件

### 将pvk格式私钥转换为pem
`openssl rsa -inform pvk -in test.pvk -outform pem -out test.pem`