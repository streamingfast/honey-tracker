{{ config(materialized='table') }}

select
    RANK() OVER(ORDER BY SUM(ap.fleet_payment) DESC) rank,
     da.address,
    sum(ap.fleet_payment) as honey,
    sum(ap.fleet_usd_payment) as USD
from hivemapper.dbt_all_payments ap
     inner join hivemapper.derived_addresses da on da.derivedaddress = ap.payee_address
group by da.address