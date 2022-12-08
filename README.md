# luxury transactions implemented with blockchain technology

# ulsan-project
2022 울산 블록체인교육 마지막 기수 미니프로젝트

# 선행조건
  - hyperledger fabric 2.2 LTS 설치
  - docker, docker-compose 설치
  - jq 설치

# 네트워크 수행
  1. 네트워크 초기화 
```
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images dev-*)
docker network prune
docker volume prune
```
  2. 네트워크 수행
```
cd 프로젝트경로/ulsan-network
./startnetwork.sh
./createchannel.sh
./setAnchorPeerUpdate.sh
./ccp-generate.sh
```
  3. 환경설정
```
export FABRIC_CFG_PATH=${PWD}/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```
  4. 채널 확인
```
peer channel list
```

# 체인코드 배포
  1. deployCC.sh 내부에 관련 정보변경 ( 체인코드 이름, 버전, 정책, 소스코드경로, 채널)
```
  ./deployCC.sh
```
  2. 체인코드 확인
```
peer lifecycle chaincode queryinstalled
peer lifecycle chaincode querycommitted -C ulsanchannel
```
# 어플리케이션 수행
```
cd 프로젝트경로/application
npm install
node server.js
```
브라우저 : 
  네트워크를 초기화 했다면 
  지갑 폴더 지우기
```
rm -rf wallet
```
  connection-org1.json 복사하기
```
  cp 프로젝트경로/ulsan-project/organizations/peerOrganization/org1.example.com/connection-org1.json ./
```

  localhost:3000 접속
  
  랜딩 페이지 접속: localhost:3000/landing.html 접속
