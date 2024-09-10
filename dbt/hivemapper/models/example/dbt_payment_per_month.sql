{{ config(materialized='table') }}

select
    DATE_TRUNC('month', p.timestamp) as month,
    sum(p.regular_payment) as regular,
    sum(p.ai_payment) as ai,
    sum(p.map_consumption_reward) as map_consumption,
    sum(p.operational_payment) as operational,
    sum(p.fleet_payment) as fleet
from hivemapper.dbt_all_payments p
group by DATE_TRUNC('month', p.timestamp)