#!/usr/bin/env bash

set -o errexit

echo "Installing Kyma..."

LOCAL_ENV=${LOCAL_ENV:-false}

CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
SCRIPTS_DIR="${CURRENT_DIR}/../scripts"
source $SCRIPTS_DIR/utils.sh
useMinikube

POSITIONAL=()
while [[ $# -gt 0 ]]
do
    key="$1"

    case ${key} in
        --kyma-release)
            checkInputParameterValue "${2}"
            KYMA_RELEASE="$2"
            shift
            shift
            ;;
         --kyma-installation)
            checkInputParameterValue "${2}"
            KYMA_INSTALLATION="$2"
            shift
            shift
            ;;
        --*)
            echo "Unknown flag ${1}"
            exit 1
        ;;
        *) # unknown option
            POSITIONAL+=("$1") # save it in an array for later
            shift # past argument
            ;;
    esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../..

INSTALLER_CR_PATH="${ROOT_PATH}"/installation/resources/kyma/installer-cr-kyma-minimal.yaml
OVERRIDES_KYMA_MINIMAL_CFG_LOCAL="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-minimal-config-local.yaml
KIALI_OVERRIDES="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kiali.yaml
CORE_TESTS_OVERRIDES="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-core-tests.yaml

INSTALLER_CR_FULL_PATH="${ROOT_PATH}"/installation/resources/kyma/installer-cr-kyma.yaml
OVERRIDES_KYMA_FULL_CFG_LOCAL="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-full-config-local.yaml

if [[ $KYMA_RELEASE == *PR-* ]]; then
  KYMA_TAG=$(curl -L https://storage.googleapis.com/kyma-development-artifacts/${KYMA_RELEASE}/kyma-installer-cluster.yaml | grep 'image: eu.gcr.io/kyma-project/kyma-installer:'| sed 's+image: eu.gcr.io/kyma-project/kyma-installer:++g' | tr -d '[:space:]')
  if [ -z "$KYMA_TAG" ]; then echo "ERROR: Kyma artifacts for ${KYMA_RELEASE} not found."; exit 1; fi
  KYMA_SOURCE="eu.gcr.io/kyma-project/kyma-installer:${KYMA_TAG}"
elif [[ $KYMA_RELEASE == master ]]; then
  KYMA_SOURCE="master"
elif [[ $KYMA_RELEASE == *master-* ]]; then
  KYMA_SOURCE=$(echo $KYMA_RELEASE | sed 's+master-++g' | tr -d '[:space:]')
else
  KYMA_SOURCE="${KYMA_RELEASE}"
fi

cert=$(</Users/i507827/.minikube/ca.crt)
cert="${cert//$'\n'/\\\\n}"

contents_minimal=$(sed "s~\"__CERT__\"~\"$cert\"~" "${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-minimal-config-local.yaml)
>"override-local-minimal.yaml" cat <<-EOF
$contents_minimal
EOF
trap "rm -f override-local-minimal.yaml" EXIT

contents_full=$(sed "s~\"__CERT__\"~\"$cert\"~" "${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-full-config-local.yaml)
>"override-local-full.yaml" cat <<-EOF
$contents_full
EOF
trap "rm -f override-local-full.yaml" EXIT

echo "Using Kyma source '${KYMA_SOURCE}'..."

echo "Installing Kyma..."
set -o xtrace
if [[ $KYMA_INSTALLATION == *full* ]]; then
  echo "Installing full Kyma"
  kyma install -c $INSTALLER_CR_FULL_PATH -o override-local-full.yaml --source $KYMA_SOURCE
else
  echo "Installing minimal Kyma"
  kyma install -c $INSTALLER_CR_PATH -o override-local-minimal.yaml --source $KYMA_SOURCE
fi
set +o xtrace