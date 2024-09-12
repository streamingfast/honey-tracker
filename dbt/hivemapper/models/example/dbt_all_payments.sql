{{
    config(
        materialized='incremental',
        indexes=[
            {'columns': ['type']},
            {'columns': ['block_number']},
            {'columns': ['block_time']},
            {'columns': ['transaction_hash']},
            {'columns': ['payee_address']},
        ]
    )
}}

select 'regular' as type,
       m.amount  as regular_payment,
       0 as ai_payment,
       0 as map_consumption_reward,
       0 as operational_payment,
       0 as fleet_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.payments p on p.mint_id = m.id
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'regular')
{% endif %}

union all

select 'ai_payments' as type,
       0  as regular_payment,
       m.amount as ai_payment,
       0 as map_consumption_reward,
       0 as operational_payment,
       0 as fleet_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.ai_payments p on p.mint_id = m.id
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'ai_payments')
{% endif %}

union all

select 'map_consumption_reward' as type,
       0  as regular_payment,
       0 as ai_payment,
       m.amount as map_consumption_reward,
       0 as operational_payment,
       0 as fleet_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.map_consumption_reward p on p.mint_id = m.id
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'map_consumption_reward')
{% endif %}

union all

select 'operational_payments' as type,
       0  as regular_payment,
       0 as ai_payment,
       0 as map_consumption_reward,
       m.amount as operational_payment,
       0 as fleet_payment,
       b.timestamp as block_time,
       b.number  as block_number,
       t.hash    as trx_hash,
       m.to_address payee_address,
       false     as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.operational_payments p on p.mint_id = m.id
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'operational_payments')
{% endif %}

union all

select 'fleet_payments' as type,
       0 as regular_payment,
       0 as ai_payment,
       0 as map_consumption_reward,
       0 as operational_payment,
       m.amount as fleet_payment,
       b.timestamp as block_time,
       b.number as block_number,
       t.hash as trx_hash,
       m.to_address as payee_address,
       p.fleet_mint_id = m.id as is_fleet
from hivemapper.blocks b
         inner join hivemapper.transactions t on t.block_id = b.id
         inner join hivemapper.mints m on m.transaction_id = t.id
         inner join hivemapper.split_payments p on p.driver_mint_id = m.id or p.fleet_mint_id = m.id
{% if is_incremental() %}
where b.number > (select coalesce(max(block_number), 0) from "hivemapper"."hivemapper"."dbt_all_payments" where type = 'fleet_payments')
{% endif %}
