
DIR="${PWD}"

echo "Installing Development Tools"
yum group install "Development Tools" -y >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "Development Tools could not be installed!" 1>&2
else
	echo "Development Tools are already installed"
fi

echo "Looking for wget"
wget --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    yum install wget zlib-devel sqlite-devel openssl openssl-devel -y
else
	echo "wget is already installed"
fi

echo "Looking for Python 3.6"
python3.6 --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	wget https://www.python.org/ftp/python/3.6.4/Python-3.6.4.tgz
 	tar -xvf Python*
	cd Python-3.6.4
	./configure
	make
	make install
else
	echo "Python 3.6 is already installed."
fi

echo "Installing required Python 3.6 Packages"
sudo pip3.6 install grpcio grpcio-tools grpcio_health_checking vcd_cli >> $DIR/setup.log 2>&1
if [ "$?" != "0" ]; then
    echo "Python packages could not be installed!" 1>&2
fi

echo "Looking for git"
git --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    yum install git
else
	echo "git is already installed"
fi

echo "Looking for go"
go version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	cd /opt
	wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
	tar -xvf go*
	export PATH=/opt/go/bin:$PATH
	export GOROOT=/opt/go
else
	echo "go is already installed"
fi

cd $DIR

echo "Installing protoc"
wget https://github.com/google/protobuf/releases/download/v3.4.1/protobuf-cpp-3.4.1.tar.gz
tar -xvf proto*
cd $DIR/proto*
yum install autoconf automake libtool curl make g++ unzip
./configure
make
make check
make install

echo "Looking for protoc"
protoc --version >> $DIR/setup.log 2>&1
if [ "$?" != "0" ]; then
    echo "[Error] protoc failed!" 1>&2
fi

cd $DIR

echo "Looking for terrform"
terraform version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	wget https://releases.hashicorp.com/terraform/0.11.4/terraform_0.11.4_darwin_amd64.zip
	unzip terraform_0.11.4_darwin_amd64.zip
	sudo mv terraform /usr/local/bin/
else
	echo "terraform is already installed"
fi

echo "Fetching Source Code"
git clone https://github.com/vmware/terraform-provider-vcloud-director.git
cd $DIR/terraform-provider-vcloud-director/go/src

export GOPATH=$DIR/terraform-provider-vcloud-director/go/
export PATH=${GOPATH}/bin:$PATH

# source ~/.bash_profile

echo "Fetching external go libraries"
./init.sh

cd $DIR

sudo mv $DIR/terraform-provider-vcloud-director/builds/linux/terraform-provider-vcloud-director /usr/local/bin/

terraform-provider-vcloud-director >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "Done Setup!"
fi
