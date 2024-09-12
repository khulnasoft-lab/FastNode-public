#!/bin/bash
(
cd ~

test -f /usr/share/fastnode/fastnode-credentials && source /usr/share/fastnode/fastnode-credentials

/usr/share/fastnode/current/fastnoded&

if [[ -z "$FASTNODE_USER" ]] || [[ -z "$FASTNODE_PASSWORD" ]]
then
      echo "FASTNODE_USER and FASTNODE_PASSWORD env variable are not defined, staying as anonymous user (Fastnode Free)"
else
      echo "Login to Fastnode with the user $FASTNODE_USER"
      sleep 5
      status_code=$( curl --write-out %{http_code} --output /usr/share/fastnode/login_logs.txt -X POST -F "email=$FASTNODE_USER" -F "password=$FASTNODE_PASSWORD" http://localhost:46624/clientapi/login)
      if [[ "$status_code" -ne 200 ]] ; then
        echo "Error while logging in, please check your credentials"
        cat /usr/share/fastnode/login_logs.txt
        echo ""
        exit 1
      else
        echo "Login to Fastnode successful with user $FASTNODE_USER"
      fi
fi

mkdir -p /home/jovyan/.local/share/fastnode/current/
cp /usr/share/fastnode/current/fastnode-lsp /home/jovyan/.local/share/fastnode/current/fastnode-lsp
cp /usr/share/fastnode/current/fastnoded /home/jovyan/.local/share/fastnode/fastnoded
) </dev/null > /usr/share/fastnode/run_logs.txt 2>&1 &

if [[ -z "$JUPYTERHUB_USER" ]]
then
  cd ~
  mkdir -p notebooks
  export NOTEBOOK_DIR_ARG="--notebook-dir=/home/$NB_USER/notebooks"
  echo "Setting notebook dir arg : $NOTEBOOK_DIR_ARG"
else
  cd
fi

echo "Command executed : start-notebook.sh $NOTEBOOK_DIR_ARG $@"
start-notebook.sh $NOTEBOOK_DIR_ARG "$@"
