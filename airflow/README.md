Fastnode Airflow
============

UI
-------------

Airflow is deploy to https://airflow.fastnode.dev. Requires VPN.

How to Deploy
-------------

Requirements:
 * AWS CLI
 * JQ (https://stedolan.github.io/jq/download/)
 * Docker

Deployment:
 * Login to AWS ECR: make docker_login
 * Deploy: make build deploy
 * Confirm Terraform deploy by type "yes"

To see deployment status:
 * make show_containers

Adding metrics to fastnode status 1d
--------------------------------

 * Ensure the field is in dags/files/fastnode_status.schema.yaml.
 * Add the aggregation to dags/templates/athena/queries/fastnode_status_1d.tmpl.sql.
 * Deploy.
 * Manually trigger the DAG "update_fastnode_status_schema": http://XXXXXXX:8080/admin/airflow/tree?dag_id=update_fastnode_status_schema
 * Let the fastnode_status_1d jobs run at their normally-scheduled time.
