`mysql`
if [ $? -ne "0" ];then
  echo "You do not have MySql instance install"
  # Have to add script for installing MySql
else
  echo "You have MySql installed, creating databses"
