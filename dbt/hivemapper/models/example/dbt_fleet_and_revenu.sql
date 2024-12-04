{{ config(materialized='table') }}

SELECT
    RANK() OVER (ORDER BY SUM(ap.fleet_payment) DESC) AS rank,
        f.address AS fleet_address,
    SUM(ap.fleet_payment) AS honey,
    SUM(ap.fleet_usd_payment) AS USD,
    COUNT(DISTINCT dp.payee_address) AS driver_count  -- Count of unique driver addresses that earned rewards
FROM
    hivemapper.fleets f
        INNER JOIN
    hivemapper.derived_addresses da ON da.address = f.address
        INNER JOIN
    hivemapper.dbt_all_payments ap ON ap.payee_address = da.derivedaddress
-- Join to identify drivers who earned rewards within the fleet
        INNER JOIN
    hivemapper.dbt_all_payments dp ON dp.trx_hash = ap.trx_hash AND dp.is_fleet = false
        INNER JOIN
    hivemapper.derived_addresses dda ON dda.derivedaddress = dp.payee_address
GROUP BY
    f.address
ORDER BY
    rank