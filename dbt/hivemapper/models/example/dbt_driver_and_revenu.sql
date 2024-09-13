{{ config(materialized='table') }}

select
    RANK() OVER(ORDER BY SUM(ap.fleet_payment) DESC) rank,
    f.address,
    sum(ap.fleet_payment) as honey,
    sum(ap.fleet_usd_payment) as USD
from hivemapper.fleets f
     inner join hivemapper.derived_addresses da on da.address = f.address
     inner join hivemapper.dbt_all_payments ap on ap.payee_address = da.derivedaddress
group by f.address