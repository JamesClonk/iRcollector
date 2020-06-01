#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [ -z "${APC_USERNAME}" ]; then
	echo "APC_USERNAME must be set!"
	exit 1
fi
if [ -z "${APC_PASSWORD}" ]; then
	echo "APC_PASSWORD must be set!"
	exit 1
fi
if [ -z "${APC_ORGANIZATION}" ]; then
	echo "APC_ORGANIZATION must be set!"
	exit 1
fi
if [ -z "${APC_SPACE}" ]; then
	echo "APC_SPACE must be set!"
	exit 1
fi
if [ -z "${LOGGLY_TOKEN}" ]; then
	echo "LOGGLY_TOKEN must be set!"
	exit 1
fi
if [ -z "${IR_USERNAME}" ]; then
	echo "IR_USERNAME must be set!"
	exit 1
fi
if [ -z "${IR_PASSWORD}" ]; then
	echo "IR_PASSWORD must be set!"
	exit 1
fi
if [ -z "${AUTH_USERNAME}" ]; then
	echo "AUTH_USERNAME must be set!"
	exit 1
fi
if [ -z "${AUTH_PASSWORD}" ]; then
	echo "AUTH_PASSWORD must be set!"
	exit 1
fi

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
echo "deploying iRcollector ..."

wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.43.0&source=github-rel' -qO cf-cli.tgz
tar -xvzf cf-cli.tgz 1>/dev/null
chmod +x cf
rm -f cf-cli.tgz || true

./cf login -a "https://api.lyra-836.appcloud.swisscom.com" -u "${APC_USERNAME}" -p "${APC_PASSWORD}" -o "${APC_ORGANIZATION}" -s "${APC_SPACE}"

# push app
./cf push iRcollector -f manifest.yml \
  --var loggly_token=${LOGGLY_TOKEN} \
  --var ir_username=${IR_USERNAME} --var ir_password=${IR_PASSWORD} \
  --var auth_username=${AUTH_USERNAME} --var auth_password=${AUTH_PASSWORD}
sleep 5

# show status
./cf app iRcollector

./cf logout

rm -f cf || true
rm -f LICENSE || true
rm -f NOTICE || true
