{{ config(materialized='table') }}

select *
from (select 'fleet' as metric, sum(fleet_payment) payment
      from hivemapper.dbt_all_payments p
      where p.type = 'fleet_payments'
        and p.is_fleet = true

      union

      select 'driver' as metric, sum(fleet_payment) payment
      from hivemapper.dbt_all_payments p
      where p.type = 'fleet_payments'
        and p.is_fleet = false) as metrics

