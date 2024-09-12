{% set execution_day = execution_date.replace(hour=0, minute=0, second=0, microsecond=0) %}
CREATE table fastnode_metrics.fastnode_status_normalized_{{ds_nodash}}
WITH (
  format='PARQUET',
  parquet_compression='SNAPPY',
  partitioned_by = ARRAY['hour'],
  external_location = 's3://fastnode-metrics/athena/fastnode_status_normalized/year={{execution_date.year}}/month={{execution_date.month}}/day={{execution_date.day}}/'
)
AS SELECT
  {% for key, value in params.schema.items()|sort if key != "properties" %}
  fastnode_status_segment.{{ key }} {{ key }},
  {% endfor %}

  {% for key, value in params.schema['properties'].items()|sort  %}
   {% if value.startswith('array') or value.startswith('map') -%}
    if(cardinality(fastnode_status_segment.properties.{{ key }}) > 0, fastnode_status_segment.properties.{{ key }}) properties__{{ key }},
   {%- else -%}
    fastnode_status_segment.properties.{{ key }} properties__{{ key }},
   {%- endif -%}
  {% endfor %}
  hour(from_iso8601_timestamp(fastnode_status_segment.timestamp)) hour
FROM fastnode_metrics.fastnode_status_segment
WHERE event IS NOT NULL
  AND event != ''
  AND prefix IN (
    '{{1000 * (execution_day - macros.timedelta(days=1)).int_timestamp}}',
    '{{1000 * execution_day.int_timestamp}}',
    '{{1000 * (execution_day + macros.timedelta(days=1)).int_timestamp}}'
  )
  AND timestamp >= '{{execution_day.strftime('%Y-%m-%dT%H:%M:%S')}}'
  AND timestamp < '{{(execution_day + macros.timedelta(days=1)).strftime('%Y-%m-%dT%H:%M:%S')}}'
