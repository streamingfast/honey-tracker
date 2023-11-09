{{ config(materialized='table') }}

select sum(mints.amount) from hivemapper.payments
inner join hivemapper.mints on mints.id = payments.mint_id