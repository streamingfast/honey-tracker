{{ config(materialized='table') }}

select DATE_TRUNC('month', b.timestamp) as month, COALESCE(SUM(br.amount), 0)  as total
from hivemapper.burns br
inner join hivemapper.transactions t on t.id = br.transaction_id
inner join hivemapper.blocks b on b.id = t.block_id
group by DATE_TRUNC('month', b.timestamp)