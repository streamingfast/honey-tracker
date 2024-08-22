{{ config(materialized='table') }}

select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
from hivemapper.payments p
inner join hivemapper.mints m on m.id = p.mint_id
inner join hivemapper.transactions t on t.id = m.transaction_id
inner join hivemapper.blocks b on b.id = t.block_id
group by DATE_TRUNC('month', b.timestamp)