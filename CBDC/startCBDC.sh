echo "*** Using network-starter ***\n"
./network-starter.sh

docker ps

source organization/nbrb/nbrb.sh
source organization/digibank/digibank.sh


echo "*** To start monitoring use: ./organization/nbrb/configuration/cli/monitordocker.sh fabric_test ***"