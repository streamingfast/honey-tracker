{{
    config(
        materialized='incremental',
        indexes=[
            {'columns': ['block_number']},
            {'columns': ['block_time']},
            {'columns': ['trx_hash']},
            {'columns': ['to_address']},
        ]
    )
}}

select b.timestamp as block_time, m.amount, m.to_address, t.hash as trx_hash, b.number as block_number
from hivemapper.mints m
         inner join hivemapper.transactions t on t.id = m.transaction_id
         inner join hivemapper.blocks b on b.id = t.block_id
{% if is_incremental() %}
where b.number > (select max(block_number) from {{ this }})
  and date_bin('5 minutes', b.timestamp , TIMESTAMP '2001-01-01') <= (select p.timestamp from "hivemapper"."hivemapper"."prices" p order by p.timestamp desc limit 1)
{% endif %}