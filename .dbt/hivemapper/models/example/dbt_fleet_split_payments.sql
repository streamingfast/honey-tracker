{{ config(materialized='table') }}

select * from (
    select 'fleet' as metric, sum(mints.amount) payment
    from hivemapper.split_payments
    inner join hivemapper.mints mints on mints.id = split_payments.fleet_mint_id

    union

    select 'driver' as metric, sum(mints.amount) payment
    from hivemapper.split_payments
    inner join hivemapper.mints mints on mints.id = split_payments.driver_mint_id
) as metrics
