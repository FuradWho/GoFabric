# NetWork

## peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name sacc --version 1.0 --sequence 1 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json

# peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name sacc --version 1.0 --init-required --package-id sacc_1:f8118a8da86e3158fc37670cca537576a4fa2e3c4c56fbb22153d2bda4515975 --sequence 1 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer lifecycle chaincode queryinstalled

# sacc_1:f8118a8da86e3158fc37670cca537576a4fa2e3c4c56fbb22153d2bda4515975

# docker cp clicpp:/opt/gopath/src/github.com/hyperledger/fabric/peer/mychannel.block ./channel-artifacts/

# peer channel join -b mychannel.block

# peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/OrgGoMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/OrgCppMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# anzhuang lianma

# docker cp sacc.tar.gz cligo:/opt/gopath/src/github.com/hyperledger/fabric/peer

# peer lifecycle chaincode install sacc.tar.gz

# sacc_1:f8118a8da86e3158fc37670cca537576a4fa2e3c4c56fbb22153d2bda4515975

# peer lifecycle chaincode queryinstalled

# peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name sacc --version 1.0 --sequence 1 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/organizations/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/organizations/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt

# peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name sacc --version 1.0 --sequence 1 --tls --cafile ${PWD}/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt

# configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate goAnchorPeer.tx -channelID mychannel -asOrg OrgGo

# peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name sacc2 --version 1.0 --sequence 2 --output json --init-required

# peer lifecycle chaincode approveformyorg --channelID mychannel --name sacc2 --version 1.0 --init-required --package-id sacc_1:f8118a8da86e3158fc37670cca537576a4fa2e3c4c56fbb22153d2bda4515975 --sequence 2 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name sacc2 --version 1.0 --sequence 1 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt --init-required

# peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/OrgCppMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# peer chaincode invoke -o orderer.example.com:7050 --isInit --ordererTLSHostnameOverride orderer.example.com --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n sacc2 --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt -c '{"Args":["a","bb"]}'

# peer chaincode query -C mychannel -n sacc2 -c '{"Args":["query","a"]}'

# peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n sacc2 --peerAddresses peer0.orggo.example.com:7051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orggo.example.com/peers/peer0.orggo.example.com/tls/ca.crt --peerAddresses peer0.orgcpp.example.com:9051 --tlsRootCertFiles ${PWD}/crypto/peerOrganizations/orgcpp.example.com/peers/peer0.orgcpp.example.com/tls/ca.crt -c '{"Args":["set","a","cc"]}'

# peer lifecycle chaincode package newchaincode.tar.gz --path /opt/gopath/src/github.com/chaincode/newchaincode/test --lang golang --label newchaincode_0

# newchaincode_0:25e24610e4e90fde66ce7debaaedebf52d4e8beb28f28b3fc72b50896dc7bd85

# peer lifecycle chaincode approveformyorg --channelID mychannel --name newchaincode_0 --version 1.0 --init-required --package-id newchaincode_0:25e24610e4e90fde66ce7debaaedebf52d4e8beb28f28b3fc72b50896dc7bd85 --sequence 1 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem