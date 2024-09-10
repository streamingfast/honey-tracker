{{ config(materialized='incremental') }}

select DATE_TRUNC('month', b.timestamp) as month, COALESCE(SUM(br.amount), 0)  as total
from hivemapper.mints br
inner join hivemapper.transactions t on t.id = br.transaction_id
inner join hivemapper.blocks b on b.id = t.block_id
{% if is_incremental() %}
    where b.timestamp >= (select coalesce(max(month),'1900-01-01') from {{ this }} )
{% endif %}

group by DATE_TRUNC('month', b.timestamp)