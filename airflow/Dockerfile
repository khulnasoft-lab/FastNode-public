FROM apache/airflow:1.10.12

USER airflow

ARG KHULNASOFT-LAB=${AIRFLOW_HOME}/khulnasoft-lab
WORKDIR ${KHULNASOFT-LAB}

COPY --chown=airflow:root airflow/requirements.txt airflow/MANIFEST.in airflow/setup.py ./airflow/
COPY --chown=airflow:root airflow/fastnode_airflow ./airflow/fastnode_airflow
COPY --chown=airflow:root fastnode-python/metrics ./fastnode-python/metrics

RUN python -m pip install --user --upgrade pip && \
    python -m pip install --user --no-cache-dir -r airflow/requirements.txt && \
    python -m pip install --user ./airflow/ && \
    python -m pip install --user ./fastnode-python/metrics/

WORKDIR ${AIRFLOW_HOME}

RUN rm -rf dags
RUN ln -s ${KHULNASOFT-LAB}/airflow/fastnode_airflow/dags dags

RUN mkdir conf
COPY airflow/conf/prod/airflow.cfg conf/
VOLUME /opt/airflow/conf
RUN ln -s conf/airflow.cfg .
