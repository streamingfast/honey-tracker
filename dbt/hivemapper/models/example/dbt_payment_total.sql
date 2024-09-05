{{ config(materialized='table') }}

with payment(total) as (
    select sum(m.amount) as total
    from hivemapper.payments p
             inner join hivemapper.mints m on m.id = p.mint_id

    union all

    select sum(m.amount) as total
    from hivemapper.ai_payments p
             inner join hivemapper.mints m on m.id = p.mint_id

    union all

    select sum(m.amount) as total
    from hivemapper.map_consumption_reward p
             inner join hivemapper.mints m on m.id = p.mint_id

    union all

    select sum(m.amount) as total
    from hivemapper.operational_payments p
             inner join hivemapper.mints m on m.id = p.mint_id

    union all

    select sum(m.amount) as total
    from hivemapper.split_payments p
             inner join hivemapper.mints m on m.id = p.fleet_mint_id or m.id = p.driver_mint_id

)
select sum(payment.total)
from payment