CREATE table fastnode_metrics.fastnode_status_normalized_{{ds_nodash}}
WITH (
  format='PARQUET',
  parquet_compression='SNAPPY',
  partitioned_by = ARRAY['hour'],
  external_location = 's3://fastnode-metrics/athena/fastnode_status_normalized/year={{execution_date.year}}/month={{execution_date.month}}/day={{execution_date.day}}/'
)
AS
WITH fastnode_status_normalized_ts AS (
  SELECT
    {% for field in params.schema['properties'] %}
    {%- if field != 'timestamp' %}{{field}},{% endif %}
    {%- endfor %}
    {#- Normalize older timestamps. Convert to ISO format and reset them based on prefix because they were client-reported and unreliable. #}
    if(regexp_like(timestamp, '^[0-9]+$'), to_iso8601(date_add('second', cast(timestamp as bigint) / 1000 - cast(to_unixtime(timestamp '{{execution_date.strftime('%Y-%m-%d %H:00')}}') as bigint), timestamp '{{execution_date.strftime('%Y-%m-%d %H:00')}}')), timestamp) timestamp
    FROM fastnode_metrics.fastnode_status
    WHERE event IS NOT NULL
      AND event != ''
      AND prefix >= '{{(execution_date.replace(hour=0, minute=0, second=0, microsecond=0) - macros.timedelta(hours=1)).strftime('%Y/%m/%d/%H')}}'
      AND prefix <= '{{(execution_date.replace(hour=0, minute=0, second=0, microsecond=0) + macros.timedelta(hours=25)).strftime('%Y/%m/%d/%H')}}'
),
fastnode_status_filtered AS (
  SELECT *,
         reduce(zip_with(split(sourceip, '.'),
                         sequence(3, 0, -1),
                         (n, p) -> cast(cast(n as bigint) * pow(256, p) as bigint)
                        ),
                cast(0 as bigint),
                (s, x) -> s + x,
                (s)->s
          ) sourceIpNumber
  FROM fastnode_status_normalized_ts
  WHERE timestamp >= '{{execution_date.replace(hour=0, minute=0, second=0, microsecond=0).strftime('%Y-%m-%dT%H:%M:%S')}}'
    AND timestamp < '{{(execution_date.replace(hour=0, minute=0, second=0, microsecond=0) + macros.timedelta(days=1)).strftime('%Y-%m-%dT%H:%M:%S')}}'
),
maxmind_masks AS (
  SELECT DISTINCT fastnode_status_filtered.sourceIp sourceip,
                  bitwise_and(fastnode_status_filtered.sourceIpNumber, maxmind.mask) maskedSourceIpNumber,
                  maxmind.mask
  FROM fastnode_status_filtered
  CROSS JOIN (SELECT DISTINCT mask FROM maxmind_city_ipv4) maxmind
  WHERE regexp_like(fastnode_status_filtered.sourceIp, '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$')
),
maxmind_cities AS (
  SELECT sourceip,
         arbitrary(maxmind.country_name) country_name,
         arbitrary(maxmind.country_iso_code) country_iso_code,
         arbitrary(maxmind.subdivision_1_name) subdivision_1_name,
         arbitrary(maxmind.city_name) city_name,
         arbitrary(maxmind.time_zone) time_zone
  FROM maxmind_masks
  JOIN maxmind_city_ipv4 maxmind
    ON maxmind_masks.mask = maxmind.mask
   AND maxmind_masks.maskedSourceIpNumber = maxmind.address
  GROUP BY sourceip
)
SELECT
  {% for key, value in params.schema['properties'].items()|sort if key != "properties" and value.get("fastnode_status_normalized") != False %}
  fastnode_status_filtered.{{ key }} {{ key }},
  {% endfor %}

  {% for key, value in params.schema['properties']['properties']['properties'].items()|sort if value.get("fastnode_status_normalized") != False %}
   {% if value.type.startswith('array') or value.type.startswith('map') -%}
    if(cardinality(fastnode_status_filtered.properties.{{ key }}) > 0, fastnode_status_filtered.properties.{{ key }}) properties__{{ key }},
   {%- else -%}
    fastnode_status_filtered.properties.{{ key }} properties__{{ key }},
   {%- endif -%}
  {% endfor %}
  maxmind_cities.country_name maxmind__country_name,
  maxmind_cities.country_iso_code maxmind__country_iso_code,
  maxmind_cities.subdivision_1_name maxmind__subdivision_1_name,
  maxmind_cities.city_name maxmind__city_name,
  maxmind_cities.time_zone maxmind__time_zone,
  monetizable_scores.score monetizable_score,
  monetizable_scores.model_version monetizable_model_version,
  hour(from_iso8601_timestamp(fastnode_status_filtered.timestamp)) hour
FROM fastnode_status_filtered
LEFT OUTER JOIN maxmind_cities
             ON fastnode_status_filtered.sourceIp = maxmind_cities.sourceip
LEFT OUTER JOIN monetizable_scores
             ON fastnode_status_filtered.userid = monetizable_scores.userid
