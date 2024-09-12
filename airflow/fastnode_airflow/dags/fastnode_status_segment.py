import datetime
import logging

from airflow import DAG
from airflow.contrib.operators.aws_athena_operator import AWSAthenaOperator
from jinja2 import PackageLoader
import fastnode_metrics
from fastnode_airflow.slack_alerts import task_fail_slack_alert


logger = logging.getLogger(__name__)


default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'start_date': datetime.datetime(2017, 4, 27),
    'end_date': datetime.datetime(2020, 2, 23),
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 0,
    'retry_delay': datetime.timedelta(minutes=5),
    'on_failure_callback': task_fail_slack_alert,
}

dag = DAG(
    'fastnode_status_segment',
    default_args=default_args,
    description='Load Segment data into fastnode_status_normalized',
    schedule_interval='10 0 * * *',
    jinja_environment_kwargs={
        'loader': PackageLoader('fastnode_airflow', 'templates')
    },
)

fastnode_status_schema = fastnode_metrics.load_schema('fastnode_status')

AWSAthenaOperator(
    aws_conn_id='aws_us_east_1',
    task_id='insert_fastnode_status_normalized',
    query='athena/queries/fastnode_status_normalized_segment.tmpl.sql',
    output_location='s3://fastnode-metrics-test/athena-results/ddl',
    database='fastnode_metrics',
    dag=dag,
    params={'schema': fastnode_status_schema}
) >> AWSAthenaOperator(
    aws_conn_id='aws_us_east_1',
    task_id='cleanup_fastnode_status_normalized_table',
    query='DROP TABLE fastnode_status_normalized_{{ds_nodash}}',
    output_location='s3://fastnode-metrics-test/athena-results/ddl',
    database='fastnode_metrics',
    dag=dag,
)
