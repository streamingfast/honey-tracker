{{ config(materialized='table') }}

select
    RANK() OVER(ORDER BY SUM(ap.fleet_payment) DESC) rank,
    da.address,
    sum(ap.fleet_payment) as honey,
    sum(ap.fleet_usd_payment) as USD
from hivemapper.dbt_all_payments ap
         inner join hivemapper.derived_addresses da on da.derivedaddress = ap.payee_address
         LEFT JOIN hivemapper.fleets f ON da.address = f.address
where f.address is null
group by da.address
