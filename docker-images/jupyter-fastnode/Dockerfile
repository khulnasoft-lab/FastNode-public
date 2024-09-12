FROM jupyter/scipy-notebook:ea01ec4d9f57
RUN pip install jupyterlab==2.2.0rc1 && pip install jupyter-fastnode && jupyter labextension install @khulnasoft-lab/jupyterlab-fastnode

USER root
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /usr/share/fastnode && mkdir /usr/share/fastnode/current
COPY fastnode-updater.sh /usr/share/fastnode
RUN chmod a+x /usr/share/fastnode/fastnode-updater.sh
RUN /usr/share/fastnode/fastnode-updater.sh --noexec --target /usr/share/fastnode/current && rm /usr/share/fastnode/fastnode-updater.sh 

RUN fix-permissions /usr/share/fastnode
USER jovyan
COPY run_fastnode.sh /usr/share/fastnode/run_fastnode.sh

ENTRYPOINT ["tini", "-g", "--"]
CMD ["/usr/share/fastnode/run_fastnode.sh"]
