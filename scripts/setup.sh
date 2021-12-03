#!/bin/bash

set -e

date "+%Y-%m-%d %H:%M:%S"
echo "------ Start executing Fabric sh script ------"

if [ "$1" = "golang" ]; then
    echo "------ Start to install golang ------"

    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt-get update
    sudo apt-get install golang-go -y
 
    sudo mkdir /home/furad/gopath
    sudo mkdir /home/furad/gopath/bin
    sudo mkdir /home/furad/gopath/pkg
    sudo mkdir -p /home/furad/gopath/src/github.com/hyperledger
 
    sudo cp /etc/profile  profile
    sudo chmod 777 profile
    echo "# GOPATH" >> profile
    echo "export GOPATH=/home/furad/gopath" >> profile 
    echo "export GOROOT=/usr/lib/go" >> profile 
    echo "export PATH=\$PATH:\$GOROOT/bin:\$GOPATH/bin " >> profile 
    sudo mv profile /etc/profile

    source /etc/profile

    echo "------ The installation is complete ! ------"

elif [ "$1" = "curl" ]; then
    echo "------ Start to install curl ------"

    sudo apt-get install libcurl3-gnutls=7.47.0-1ubuntu2.19
    sudo apt-get install curl

    echo "------ The installation is complete ! ------"

elif [ "$1" = "docker" ]; then
    echo "------ Start to install docker ------"

    # 安装 docker 和 docker-compose
    sudo apt-get remove docker docker-engine docker.io containerd runc
    sudo apt-get update

    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) \
    stable"
    sudo apt-get update
    sudo apt-get install docker-ce docker-ce-cli containerd.io -y

    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sleep 5
    sudo chmod +x /usr/local/bin/docker-compose
    
    sudo groupadd docker
    sudo gpasswd -a $USER docker
    newgrp docker
 
    sudo mkdir -p /etc/docker
    echo "{" > daemon.json
    echo "\"registry-mirrors\": [\"https://ofyci9nf.mirror.aliyuncs.com\"]" >> daemon.json
    echo "}" >> daemon.json
    echo "" >> daemon.json
 
    sudo mv daemon.json  /etc/docker/daemon.json
 
    sudo systemctl daemon-reload     
    sudo systemctl restart docker

    echo "------ The installation is complete ! ------"

elif [ "$1" = "fabric" ]; then
    echo "------ Start to install fabric ------"

#    sudo cp /etc/hosts  hosts
#    sudo chmod 777 hosts
#    echo "199.232.68.133 raw.githubusercontent.com" >> hosts
#    echo "199.232.68.133 user-images.githubusercontent.com" >> hosts
#    echo "199.232.68.133 avatars2.githubusercontent.com" >> hosts
#    echo "199.232.68.133 avatars1.githubusercontent.com" >> hosts
#    sudo mv hosts /etc/hosts

    sudo mkdir -p $GOPATH/src/github.com/hyperledger
    cd $GOPATH/src/github.com/hyperledger
    sudo rm -rf  $GOPATH/src/github.com/hyperledger/fabric-samples # 如果有的话先删掉
    # sudo curl -sS https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh -o bootsrap.sh
    sudo chmod 777 bootstrap.sh
    sudo ./bootstrap.sh
    sudo chmod 777 -R fabric-samples

    sudo cp /etc/profile profile
    sudo chmod 777 profile
    echo "# fabric-samples-bin" >> profile
    echo "export PATH=\$PATH:\$GOPATH/src/github.com/hyperledger/fabric-samples/bin" >> profile 
 
    sudo mv profile /etc/profile
    source /etc/profile

    echo "------ The installation is complete ! ------"

elif [ "$1" = "git" ]; then
    echo "------ Start to install git ------"

    sudo apt-get install git

    echo "------ The installation is complete ! ------"

elif [ "$1" = "update" ]; then
    echo "------ Start to install update ------"
    sudo cp /etc/apt/sources.list /etc/apt/sources.list.bak

    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse" > sources.list
    echo "deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse" >> sources.list
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse" >> sources.list
    echo "deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse" >> sources.list
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse" >> sources.list
    echo "deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse" >> sources.list
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse" >> sources.list
    echo "deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse" >> sources.list
    sudo mv sources.list /etc/apt/sources.list
    sudo apt-get update

    echo "------ The installation is complete ! ------"

elif [ "$1" = "help" ]; then
    
    echo "------ setup.sh shell help ------"
    echo -e "\n"
    echo "------ Usage: bash setup.sh COMMAND ------"
    echo -e "\n"
    echo "COMMANDS: "
    echo "    update    Update the apt-get"
    echo "    git       Start to install git"
    echo "    curl      Start to install curl"
    echo "    golang    Start to install golang"
    echo "    docker    Start to install docker"
    echo "    fabric    Start to install fabric"
    echo -e "\n"

fi