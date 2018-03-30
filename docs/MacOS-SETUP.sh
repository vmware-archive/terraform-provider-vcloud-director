
DIR="${PWD}"

echo "Looking for xcode-select"
xcode-select --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    xcode-select --install
else
	echo "xcode-select is already installed"
fi

echo "Looking for brew"
brew --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
else
	echo "brew is already installed"
fi


echo "Looking for wget"
wget --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    brew install wget --with-libressl
else
	echo "wget is already installed"
fi

echo "Looking for Python 3.6"
python3 --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    brew install python
else
	echo "Python 3.6 is already installed"
fi

echo "Installing required Python Packages"
sudo pip3.6 install grpcio grpcio-tools grpcio_health_checking vcd_cli >> $DIR/setup.log 2>&1
if [ "$?" != "0" ]; then
    echo "Python packages could not be installed!" 1>&2
fi

echo "Looking for git"
git --version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    brew install git
else
	echo "git is already installed"
fi

echo "Looking for go"
go version >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
    brew install go
else
	echo "go is already installed"
fi

export GOVERSION=$(brew list go | head -n 1 | cut -d '/' -f 6)
export GOROOT=$(brew --prefix)/Cellar/go/${GOVERSION}/libexec
# source ~/.bash_profile

echo "Installing protoc"
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

sudo mv $DIR/terraform-provider-vcloud-director/builds/mac/terraform-provider-vcloud-director /usr/local/bin/

terraform-provider-vcloud-director >> $DIR/setup.log 2>&1
if [ "$?" -ne 0 ]; then
	echo "Done Setup!"
fi


