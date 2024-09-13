{{ config(materialized='table') }}

select
    DATE_TRUNC('month', p.block_time) as month,
    sum(p.amount) as total
from hivemapper.dbt_all_mints p
group by DATE_TRUNC('month', p.block_time)