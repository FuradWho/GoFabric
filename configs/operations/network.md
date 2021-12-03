# NetWork 操作



**创建通道**
```shell
peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

**加入通道**
```shell
peer channel join -b mychannel.block
```

**显示已加入的通道**
```shell
peer channel list
```

**docker 与 宿主机复制文件**
```shell
docker cp container:xxx  xxx
```

**更新锚节点**
```shell
peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/OrgGoMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

**安装链码**
```shell
peer lifecycle chaincode install sacc.tar.gz
```

**查询安装的链码**
```shell
peer lifecycle chaincode queryinstalled
```

**申请当前所在组织批准链码**
```shell
peer lifecycle chaincode approveformyorg --channelID mychannel --name newchaincode_0 --version 1.0 --init-required --package-id newchaincode_0:25e24610e4e90fde66ce7debaaedebf52d4e8beb28f28b3fc72b50896dc7bd85 --sequence 1 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

**查询通道内组织批准详情**
```shell
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name sacc2 --version 1.0 --sequence 2 --output json --init-required
```

**提交通过批准的链码**
```shell
peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name sacc2 --version 1.0 --sequence 1 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt --init-required
```
