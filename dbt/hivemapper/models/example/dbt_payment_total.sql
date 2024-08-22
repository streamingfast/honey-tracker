{{ config(materialized='table') }}

select sum(mints.amount) as total_payments
from hivemapper.payments
inner join hivemapper.mints on mints.id = payments.mint_id
