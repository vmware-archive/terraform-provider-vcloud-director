#!/bin/bash

DIR="${PWD}"

echo "CHECK - Looking for xcode-select"
xcode-select --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing xcode-select" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing xcode-select"
    xcode-select --install
else
	echo "DONE - xcode-select is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - xcode-select is already installed"
fi

echo "CHECK - Looking for brew"
brew --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing brew" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing brew"
    ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
else
	echo "DONE - brew is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - brew is already installed"
fi


echo "CHECK - Looking for wget"
wget --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing wget" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing wget"
    brew install wget --with-libressl
else
	echo "DONE - wget is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - wget is already installed"
fi

echo "CHECK - Looking for Python 3"
python3 --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing Python 3" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing Python 3"
    brew install python
else
	echo "DONE - Python 3 is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - Python 3 is already installed"
fi

echo "CHECK - Looking for required Python Packages"
pip3.6 install grpcio grpcio-tools grpcio_health_checking vcd_cli >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    echo "ERROR - Python packages could not be installed!" 1>&2
else
	echo "DONE - Installed Python Packages"
fi

echo "CHECK - Looking for git"
git --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing git" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing git"
    brew install git
else
	echo "DONE - git is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - git is already installed"
fi

echo "CHECK - Looking for go"
go version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing go" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing go"
    brew install go
else
	echo "DONE - go is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - go is already installed"
fi

export GOVERSION=$(brew list go | head -n 1 | cut -d '/' -f 6)
export GOROOT=$(brew --prefix)/Cellar/go/${GOVERSION}/libexec
# source ~/.bash_profile

echo "CHECK - Looking for protoc"
protoc --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing protoc" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing protoc"
	wget https://github.com/google/protobuf/releases/download/v3.4.1/protobuf-cpp-3.4.1.tar.gz
	tar -xvf proto*
	cd $DIR/proto*
	wget https://gist.githubusercontent.com/justinbellamy/2672db1c78f024f2d4fe/raw/617e39f18f32a6a97a365dedafdc93137c625738/cltools.sh
	chmod +x cltools.sh
	./cltools.sh
	./configure
	make
	make check
	make install
	cd $DIR
	rm -rf protobuf*
else
	echo "DONE - protoc is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - protoc is already installed"
fi

echo "CHECK - Looking for terrform"
terraform version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "INSTALL - Installing Terraform" >> $DIR/setup.log 2>&1
	echo "INSTALL - Installing Terraform"
	wget https://releases.hashicorp.com/terraform/0.11.4/terraform_0.11.4_darwin_amd64.zip
	unzip terraform_0.11.4_darwin_amd64.zip
	sudo mv terraform /usr/local/bin/
else
	echo "DONE - Terraform is already installed" >> $DIR/setup.log 2>&1
	echo "DONE - Terraform is already installed"
fi


function fetch_source_code(){
	if [ $1 == "fetch" ]; then
		echo "FETCH - Fetching latest Terraform Provider Plugin source code" >> $DIR/setup.log 2>&1
		echo "FETCH - Fetching latest Terraform Provider Plugin source code"
		git clone https://github.com/vmware/terraform-provider-vcloud-director.git
	fi
	cd $DIR/terraform-provider-vcloud-director/go/src
	export PROJECTDIR=$DIR/terraform-provider-vcloud-director
	export GOPATH=$DIR/terraform-provider-vcloud-director/go/
	export PATH=${GOPATH}/bin:$PATH
	echo "FETCH - Fetching external go libraries" >> $DIR/setup.log 2>&1
	echo "FETCH - Fetching external go libraries"
	./init.sh
}

if [ -d "terraform-provider-vcloud-director" ]; then
	echo -e "INFO - Terraform Provider Plugin source directory is already exists"
	echo -e "Do you want to fetch the latest source code (Y/N)? \c"
	read input
	if [ $input == "Y" ]; then
		rm -rf terraform-provider-vcloud-director
		fetch_source_code fetch
	else
		fetch_source_code build
	fi
else
	fetch_source_code fetch
fi

cd $DIR
sudo mv $DIR/terraform-provider-vcloud-director/builds/mac/terraform-provider-vcloud-director /usr/local/bin/

terraform-provider-vcloud-director >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "Done Setup!" >> $DIR/setup.log 2>&1
	echo "Done Setup!"
fi


