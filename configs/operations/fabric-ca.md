// login tls ca server

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/tls-ca/admin fabric-ca-client enroll -d
-u https://tls-ca-admin:tls-ca-adminpw@tls-ca.example.com:7052

// Register TLS certificate for each component (the component refers to peer,order and administrator). It only registers
the identity, but does not obtain the certificate;

fabric-ca-client register -d --id.name peer1-org1 --id.secret peer1PW --id.type peer -u https://tls-ca.example.com:7052
fabric-ca-client register -d --id.name peer2-org1 --id.secret peer2PW --id.type peer -u https://tls-ca.example.com:7052
fabric-ca-client register -d --id.name peer1-org2 --id.secret peer1PW --id.type peer -u https://tls-ca.example.com:7052
fabric-ca-client register -d --id.name peer2-org2 --id.secret peer2PW --id.type peer -u https://tls-ca.example.com:7052

// admin administrator registered with org1

fabric-ca-client register -d --id.name admin-org1 --id.secret org1AdminPW --id.type admin
-u https://tls-ca.example.com:7052

// Register the admin administrator of org2

fabric-ca-client register -d --id.name admin-org2 --id.secret org2AdminPW --id.type admin
-u https://tls-ca.example.com:7052

fabric-ca-client register -d --id.name orderer-org0 --id.secret ordererPW --id.type orderer
-u https://tls-ca.example.com:7052

//
-------------------------------------------------------------------------------------------------------------------------------------------------

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org0-ca/admin

fabric-ca-client enroll -d -u https://ca-org0-admin:ca-org0-adminpw@ca.org0.example.com:7053

// Register orderer user fabric-ca-client register -d --id.name orderer-org0 --id.secret ordererPW --id.type orderer
--id.attrs '"hf.Registrar.Roles=orderer"' -u https://ca.org0.example.com:7053

// Register admin user fabric-ca-client register -d --id.name admin2-org0 --id.secret org0Admin2Pw --id.type client
--id.attrs "hf.Registrar.Roles=client,orderer,peer,user,hf.Registrar.Attributes=*
,hf.Revoker=true,hf.GenCRL=true,role=admin=true:ecert" -u https://ca.org0.example.com:7053

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/orderer-org0 export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://orderer-org0:ordererPW@ca.org0.example.com:7053

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/orderer-org0 export FABRIC_CA_CLIENT_MSPDIR=tls-msp
fabric-ca-client enroll -d -u https://orderer-org0:ordererPW@tls-ca.example.com:7052 --enrollment.profile tls
--csr.hosts orderer.org0.example.com

mv /usr/local/hyper/test3/crypto-config/orderer-org0/tls-msp/keystore/*_sk
/usr/local/hyper/test3/crypto-config/orderer-org0/tls-msp/keystore/key.pem

*/

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/orderer-org0-admin export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin-org0:org0AdminPw@ca.org0.example.com:7053

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/orderer-org0-admin1 export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin1-org0:org0Admin1Pw@ca.org0.example.com:7053

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/orderer-org0-admin2 export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin2-org0:org0Admin2Pw@ca.org0.example.com:7053

mkdir /usr/local/hyper/test3/crypto-config/orderer-org0/msp/admincerts cp
/usr/local/hyper/test3/crypto-config/orderer-org0-admin/msp/signcerts/cert.pem
/usr/local/hyper/test3/crypto-config/orderer-org0/msp/admincerts

//
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org1-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-ca/admin export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://ca-org1-admin:ca-org1-adminpw@ca.org1.example.com:7054

export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-ca/admin export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://ca-org1-admin:ca-org1-adminpw@ca.org1.example.com:7054

export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test5/fabric-ca/org1-ca/admin export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u http://ca-org1-admin:ca-org1-adminpw@ca.org1.example.com:7054

fabric-ca-client register -d --id.name peer1-org1 --id.secret peer1PW --id.type peer -u https://ca.org1.example.com:7054
fabric-ca-client register -d --id.name peer2-org1 --id.secret peer2PW --id.type peer -u https://ca.org1.example.com:7054
fabric-ca-client register -d --id.name admin-org1 --id.secret org1AdminPW --id.type admin
-u https://ca.org1.example.com:7054

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org1-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-admin
fabric-ca-client enroll -d -u https://admin-org1:org1AdminPW@ca.org1.example.com:7054

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org1-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-peer1
fabric-ca-client enroll -d -u https://peer1-org1:peer1PW@ca.org1.example.com:7054

# 拷贝证书

mkdir /usr/local/hyper/test3/crypto-config/org1-peer1/msp/admincerts cp
/usr/local/hyper/test3/crypto-config/org1-admin/msp/signcerts/cert.pem
/usr/local/hyper/test3/crypto-config/org1-peer1/msp/admincerts

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=tls-msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-peer1
fabric-ca-client enroll -d -u https://peer1-org1:peer1PW@tls-ca.example.com:7052 --enrollment.profile tls --csr.hosts
peer1.org1.example.com

# 重命名

mv /usr/local/hyper/test3/crypto-config/org1-peer1/tls-msp/keystore/*_sk
/usr/local/hyper/test3/crypto-config/org1-peer1/tls-msp/keystore/key.pem
*/

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org1-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-peer2
fabric-ca-client enroll -d -u https://peer2-org1:peer2PW@ca.org1.example.com:7054

# 拷贝证书

mkdir /usr/local/hyper/test3/crypto-config/org1-peer2/msp/admincerts cp
/usr/local/hyper/test3/crypto-config/org1-admin/msp/signcerts/cert.pem
/usr/local/hyper/test3/crypto-config/org1-peer2/msp/admincerts

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=tls-msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org1-peer2
fabric-ca-client enroll -d -u https://peer2-org1:peer2PW@tls-ca.example.com:7052 --enrollment.profile tls --csr.hosts
peer2.org1.example.com

# 重命名

mv /usr/local/hyper/test3/crypto-config/org1-peer2/tls-msp/keystore/*_sk
/usr/local/hyper/test3/crypto-config/org1-peer2/tls-msp/keystore/key.pem
*/

//
--------------------------------------------------------------------------------------------------------------------------------------------------------------

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org2-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-ca/admin export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://ca-org2-admin:ca-org2-adminpw@ca.org2.example.com:7055

fabric-ca-client register -d --id.name peer1-org2 --id.secret peer1PW --id.type peer -u https://ca.org2.example.com:7055
fabric-ca-client register -d --id.name peer2-org2 --id.secret peer2PW --id.type peer -u https://ca.org2.example.com:7055
fabric-ca-client register -d --id.name admin-org2 --id.secret org2AdminPW --id.type admin
-u https://ca.org2.example.com:7055

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org2-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-admin
fabric-ca-client enroll -d -u https://admin-org2:org2AdminPW@ca.org2.example.com:7055

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org2-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-peer1
fabric-ca-client enroll -d -u https://peer1-org2:peer1PW@ca.org2.example.com:7055

mkdir /usr/local/hyper/test3/crypto-config/org2-peer1/msp/admincerts cp
/usr/local/hyper/test3/crypto-config/org2-admin/msp/signcerts/cert.pem
/usr/local/hyper/test3/crypto-config/org2-peer1/msp/admincerts

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=tls-msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-peer1
fabric-ca-client enroll -d -u https://peer1-org2:peer1PW@tls-ca.example.com:7052 --enrollment.profile tls --csr.hosts
peer1.org2.example.com

mv /usr/local/hyper/test3/crypto-config/org2-peer1/tls-msp/keystore/*_sk
/usr/local/hyper/test3/crypto-config/org2-peer1/tls-msp/keystore/key.pem

*/

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/org2-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-peer2
fabric-ca-client enroll -d -u https://peer2-org2:peer2PW@ca.org2.example.com:7055

mkdir /usr/local/hyper/test3/crypto-config/org2-peer2/msp/admincerts cp
/usr/local/hyper/test3/crypto-config/org2-admin/msp/signcerts/cert.pem
/usr/local/hyper/test3/crypto-config/org2-peer2/msp/admincerts

export FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem export
FABRIC_CA_CLIENT_MSPDIR=tls-msp export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/test3/crypto-config/org2-peer2
fabric-ca-client enroll -d -u https://peer2-org2:peer2PW@tls-ca.example.com:7052 --enrollment.profile tls --csr.hosts
peer2.org2.example.com

mv /usr/local/hyper/test3/crypto-config/org2-peer2/tls-msp/keystore/*_sk
/usr/local/hyper/test3/crypto-config/org2-peer2/tls-msp/keystore/key.pem
*/

//
-------------------------------------------------------------------------------------------------------------------------------------------------------------------

mkdir msp/admincerts

cp /usr/local/hyper/test3/crypto-config/orderer-org0-admin/msp/signcerts/cert.pem msp/admincerts/cert.pem

mkdir msp/cacerts

cp /usr/local/hyper/test3/crypto-config/org0-ca/crypto/ca-cert.pem msp/cacerts/cert.pem

mkdir msp/tlscacerts cp /usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem msp/tlscacerts/cert.pem

mkdir msp/admincerts

cp /usr/local/hyper/test3/crypto-config/org1-admin/msp/signcerts/cert.pem msp/admincerts/cert.pem

mkdir msp/cacerts

cp /usr/local/hyper/test3/crypto-config/org1-ca/crypto/ca-cert.pem msp/cacerts/cert.pem

mkdir msp/tlscacerts cp /usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem msp/tlscacerts/cert.pem

mkdir msp/admincerts

cp /usr/local/hyper/test3/crypto-config/org2-admin/msp/signcerts/cert.pem msp/admincerts/cert.pem

mkdir msp/cacerts

cp /usr/local/hyper/test3/crypto-config/org2-ca/crypto/ca-cert.pem msp/cacerts/cert.pem

mkdir msp/tlscacerts cp /usr/local/hyper/test3/crypto-config/tls-ca/crypto/ca-cert.pem msp/tlscacerts/cert.pem

// Use the super administrator account and password org1-admin:org1-adminpw logon server fabric-ca-client enroll -d
-u https://rca-org1-admin:rca-org1-adminpw@0.0.0.0:17054 --tls.certfiles
/usr/local/hyper/fabric-ca/org1/ca/crypto/ca-cert.pem

// Note the change in the type value fabric-ca-client register -d --id.name peer1-org1 --id.secret peer1PW --id.type
peer -u https://0.0.0.0:17054
fabric-ca-client register -d --id.name peer2-org1 --id.secret peer2PW --id.type peer -u https://0.0.0.0:17054
fabric-ca-client register -d --id.name admin-org1 --id.secret org1AdminPW --id.type admin -u https://0.0.0.0:17054

// /usr/local/home/org1/ca/crypto/ca-cert.pem Is the root certificate of org1 organization generated after the service
is started. This certificate is required for login export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org2/ca/crypto/ca-cert.pem

// After login, super administrator root certificate will be generated under / usr/local/home/org0/ca/admin export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org2/ca/admin

// Use the super administrator account and password org1-admin:org1-adminpw logon server fabric-ca-client enroll -d
-u https://rca-org2-admin:rca-org2-adminpw@0.0.0.0:17055 --tls.certfiles
/usr/local/hyper/fabric-ca/org2/ca/crypto/ca-cert.pem

// Register all nodes in org2 organization, including peer1, peer2, and Admin. Note that the type values of admin and
peer are different fabric-ca-client register -d --id.name peer1-org2 --id.secret peer1PW --id.type peer
-u https://0.0.0.0:17055
fabric-ca-client register -d --id.name peer2-org2 --id.secret peer2PW --id.type peer -u https://0.0.0.0:17055
fabric-ca-client register -d --id.name admin-org2 --id.secret org2AdminPW --id.type admin -u https://0.0.0.0:17055

// peer1-org1 stores the root directory of the certificate export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org1/peer1

// The tls certificate communicating with org1 CA uses the certificate generated when org1 service is started export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org1/ca/crypto/ca-cert.pem

// Directory used to store msp certificates for peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://peer1-org1:peer1PW@0.0.0.0:17054

// Save tls certificate directory of peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Using root certificate of TLS CA to access TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://peer1-org1:peer1PW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts peer1-org1

mv /usr/local/hyper/fabric-ca/org1/peer1/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org1/peer1/tls-msp/keystore/key.pem

*/

// peer1-org1 stores the root directory of the certificate export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org1/peer2

// The tls certificate communicating with org1 CA uses the certificate generated when org1 service is started export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org1/ca/crypto/ca-cert.pem

// Directory used to store msp certificates for peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://peer2-org1:peer2PW@0.0.0.0:17054

// Save tls certificate directory of peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Using root certificate of TLS CA to access TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://peer2-org1:peer2PW@0.0.0.0:17052

mv /usr/local/hyper/fabric-ca/org1/peer2/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org1/peer2/tls-msp/keystore/key.pem

// peer1-org1 stores the root directory of the certificate export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org1/admin

// The tls certificate communicating with org1 CA uses the certificate generated when org1 service is started export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org1/ca/crypto/ca-cert.pem

// Directory used to store msp certificates for peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://admin-org1:org1AdminPW@0.0.0.0:17054

// Save tls certificate directory of peer1-org1 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Using root certificate of TLS CA to access TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://admin-org1:org1AdminPW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts
admin-org1

mv /usr/local/hyper/fabric-ca/org1/admin/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org1/admin/tls-msp/keystore/key.pem

mkdir -p /usr/local/hyper/fabric-ca/org1/peer1/msp/admincerts && cp
/usr/local/hyper/fabric-ca/org1/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org1/peer1/msp/admincerts/org1-admin-cert.pem

mkdir -p /usr/local/hyper/fabric-ca/org1/peer2/msp/admincerts && cp
/usr/local/hyper/fabric-ca/org1/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org1/peer2/msp/admincerts/org1-admin-cert.pem

// Set the certificate root directory of peer1-org2 export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org2/peer1

// Using root certificate of org2 CA server to communicate with org2 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org2/ca/crypto/ca-cert.pem

// Set the certificate directory of msp of peer1-org2 export FABRIC_CA_CLIENT_MSPDIR=msp

// Log in to org2 CA server with peer1-org2 account fabric-ca-client enroll -d -u https://peer1-org2:peer1PW@0.0.0.0:
17055

// Set tls certificate directory of peer1-org2 export FABRIC_CA_CLIENT_MSPDIR=tls-msp // Using root certificate of TLS
CA server to communicate with TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://peer1-org2:peer1PW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts peer1-org2

mv /usr/local/hyper/fabric-ca/org2/peer1/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org2/peer1/tls-msp/keystore/key.pem

// Set the certificate root directory of peer2-org2 export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org2/peer2

// Using root certificate of org2 CA server to communicate with org2 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org2/ca/crypto/ca-cert.pem

// Set the certificate directory of msp of peer2-org2 export FABRIC_CA_CLIENT_MSPDIR=msp

// Log in to org2 CA server with peer2-org2 account fabric-ca-client enroll -d -u https://peer2-org2:peer2PW@0.0.0.0:
17055

After setting up the environment of peer2-org2, log in to tls CA server with peer2-org2, and obtain tls certificate of
peer2-org2

// Set tls certificate directory of peer2-org2 export FABRIC_CA_CLIENT_MSPDIR=tls-msp // Using root certificate of TLS
CA server to communicate with TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://peer2-org2:peer2PW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts peer2-org2

mv /usr/local/hyper/fabric-ca/org2/peer2/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org2/peer2/tls-msp/keystore/key.pem

// Set the certificate root directory of admin-org2 export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org2/admin

// Using root certificate of org2 CA server to communicate with org2 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org2/ca/crypto/ca-cert.pem

// Set the certificate directory of msp of admin-org2 export FABRIC_CA_CLIENT_MSPDIR=msp

// Log in to org2 CA server with admin-org2 account fabric-ca-client enroll -d -u https://admin-org2:
org2AdminPW@0.0.0.0:17055

// Set tls certificate directory of admin-org2 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Using root certificate of TLS CA server to communicate with TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://admin-org2:org2AdminPW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts
admin-org2

mv /usr/local/hyper/fabric-ca/org2/admin/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org2/admin/tls-msp/keystore/key.pem

mkdir -p /usr/local/hyper/fabric-ca/org2/peer1/msp/admincerts && cp
/usr/local/hyper/fabric-ca/org2/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org2/peer1/msp/admincerts/org2-admin-cert.pem

mkdir -p /usr/local/hyper/fabric-ca/org2/peer2/msp/admincerts && cp
/usr/local/hyper/fabric-ca/org2/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org2/peer2/msp/admincerts/org2-admin-cert.pem

mkdir -p /usr/local/hyper/fabric-ca/org0/orderers

// Set the msp certificate root directory of orderer1-org0 export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0

// Use the certificate of org0 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org0/ca/crypto/ca-cert.pem

// Set the msp certificate directory of orderer1-org0 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://orderer1-org0:ordererpw@ca.org0.example.com:7053

// Set the tls certificate root directory of orderer1-org0 export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0

// Set tls certificate directory of orderer1-org0 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Use the certificate generated when TLS CA server starts to communicate with TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer1-org0:ordererPW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts
orderer1-org0

mv /usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0/tls-msp/keystore/key.pem

// Set the certificate root directory of admin-org0 export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/admin

// Use the certificate generated when org0 CA server starts to communicate with org0 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org0/ca/crypto/ca-cert.pem

// Set the msp certificate directory of admin-org0 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://admin-org0:org0adminpw@ca.org0.example.com:7053

mkdir /usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0/msp/admincerts cp
/usr/local/hyper/fabric-ca/org0/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org0/orderers/orderer1-org0/msp/admincerts/orderer-admin-cert.pem

// Set the msp certificate root directory of orderer1-org0 export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0

// Use the certificate of org0 CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org0/ca/crypto/ca-cert.pem

// Set the msp certificate directory of orderer1-org0 export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://orderer2-org0:ordererpw@ca.org0.example.com:7053

// Set the tls certificate root directory of orderer1-org0 export
FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0

// Set tls certificate directory of orderer1-org0 export FABRIC_CA_CLIENT_MSPDIR=tls-msp

// Use the certificate generated when TLS CA server starts to communicate with TLS CA server export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer2-org0:ordererPW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts
orderer2-org0

mv /usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0/tls-msp/keystore/key.pem

mkdir /usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0/msp/admincerts cp
/usr/local/hyper/fabric-ca/org0/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org0/orderers/orderer2-org0/msp/admincerts/orderer-admin-cert.pem

export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0 export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/org0/ca/crypto/ca-cert.pem export FABRIC_CA_CLIENT_MSPDIR=msp

fabric-ca-client enroll -d -u https://orderer3-org0:ordererpw@ca.org0.example.com:7053

export FABRIC_CA_CLIENT_HOME=/usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0 export
FABRIC_CA_CLIENT_MSPDIR=tls-msp export
FABRIC_CA_CLIENT_TLS_CERTFILES=/usr/local/hyper/fabric-ca/tls-ca/crypto/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer3-org0:ordererPW@0.0.0.0:17052 --enrollment.profile tls --csr.hosts
orderer3-org0

mv /usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0/tls-msp/keystore/*_sk
/usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0/tls-msp/keystore/key.pem

mkdir /usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0/msp/admincerts cp
/usr/local/hyper/fabric-ca/org0/admin/msp/signcerts/cert.pem
/usr/local/hyper/fabric-ca/org0/orderers/orderer3-org0/msp/admincerts/orderer-admin-cert.pem

mkdir -p /usr/local/hyper/fabric-ca/configtx/org1/msp/tlscacerts && cp
/usr/local/hyper/fabric-ca/org1/admin/tls-msp/tlscacerts/tls-0-0-0-0-17052.pem
/usr/local/hyper/fabric-ca/configtx/org1/tlscacerts

mkdir -p /usr/local/hyper/fabric-ca/configtx/{block,channel-artifacts}

configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock
/usr/local/hyper/fabric-ca/configtx/block/genesis.block

export CHANNEL_NAME=mychannel

configtxgen -profile TwoOrgsChannel -outputCreateChannelTx
/usr/local/hyper/fabric-ca/configtx/channel-artifacts/${CHANNEL_NAME}.tx -channelID ${CHANNEL_NAME}

export orgmsp=org1MSP configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate
/usr/local/hyper/fabric-ca/configtx/channel-artifacts/${orgmsp}anchors.tx -channelID ${CHANNEL_NAME} -asOrg ${orgmsp}

export orgmsp=org2MSP configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate
/usr/local/hyper/fabric-ca/configtx/channel-artifacts/${orgmsp}anchors.tx -channelID ${CHANNEL_NAME} -asOrg ${orgmsp}


