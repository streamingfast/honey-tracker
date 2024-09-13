{{
    config(
        materialized='incremental',
        indexes=[
            {'columns': ['type']},
            {'columns': ['block_number']},
            {'columns': ['block_time']},
            {'columns': ['trx_hash']},
            {'columns': ['payee_address']},
        ]
    )
}}

select 'regular' as type,
       m.amount  as regular_payment,
       prices.price * m.amount as regular_usd_payment,
       0 as ai_payment,
       0 as ai_usd_payment,
       0 as map_consumption_reward,
       0 as map_consumption_usd_reward,
       0 as operational_payment,
       0 as operational_usd_payment,
       0 as fleet_payment,
       0 as fleet_usd_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.payments p on p.mint_id = m.id
         inner join hivemapper.prices prices on prices.timestamp = date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01')
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'regular')
and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}

union all

select 'ai_payments' as type,
    0  as regular_payment,
    0 as regular_usd_payment,
    m.amount as ai_payment,
    prices.price * m.amount as ai_usd_payment,
    0 as map_consumption_reward,
    0 as map_consumption_usd_reward,
    0 as operational_payment,
    0 as operational_usd_payment,
    0 as fleet_payment,
    0 as fleet_usd_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.ai_payments p on p.mint_id = m.id
    inner join hivemapper.prices prices on prices.timestamp = date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01')
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'ai_payments')
and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}

union all

select 'map_consumption_reward' as type,
    0  as regular_payment,
    0 as regular_usd_payment,
    0 as ai_payment,
    0 as ai_usd_payment,
    m.amount as map_consumption_reward,
    prices.price * m.amount as map_consumption_usd_reward,
    0 as operational_payment,
    0 as operational_usd_payment,
    0 as fleet_payment,
    0 as fleet_usd_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.map_consumption_reward p on p.mint_id = m.id
         inner join hivemapper.prices prices on prices.timestamp = date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01')
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'map_consumption_reward')
and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}

union all

select 'operational_payments' as type,
    0  as regular_payment,
    0 as regular_usd_payment,
    0 as ai_payment,
    0 as ai_usd_payment,
    0 as map_consumption_reward,
    0 as map_consumption_usd_reward,
    m.amount as operational_payment,
    prices.price * m.amount as operational_usd_payment,
    0 as fleet_payment,
    0 as fleet_usd_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.operational_payments p on p.mint_id = m.id
    inner join hivemapper.prices prices on prices.timestamp = date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01')
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'operational_payments')
  and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}

union all

select 'fleet_payments' as type,
    0  as regular_payment,
    0 as regular_usd_payment,
    0 as ai_payment,
    0 as ai_usd_payment,
    0 as map_consumption_reward,
    0 as map_consumption_usd_reward,
    0 as operational_payment,
    0 as operational_usd_payment,
    m.amount as fleet_payment,
    prices.price * m.amount as fleet_usd_payment,
       b.timestamp as block_time,
       b.number as block_number,
       t.hash as trx_hash,
       m.to_address as payee_address,
       p.fleet_mint_id = m.id as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.split_payments p on p.driver_mint_id = m.id or p.fleet_mint_id = m.id
    inner join hivemapper.prices prices on prices.timestamp = date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01')
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'fleet_payments')
  and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}
