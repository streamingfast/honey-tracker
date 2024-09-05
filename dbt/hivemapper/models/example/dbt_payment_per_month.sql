{{ config(materialized='table') }}

-- select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
-- from hivemapper.payments p
-- inner join hivemapper.mints m on m.id = p.mint_id
-- inner join hivemapper.transactions t on t.id = m.transaction_id
-- inner join hivemapper.blocks b on b.id = t.block_id
-- group by DATE_TRUNC('month', b.timestamp)


select
    payment.month as month,
    COALESCE(payment.total,0) driver_payment,
    COALESCE(ai_payment.total,0) as ai_payment,
    COALESCE(map_consumption_reward.total,0) as map_consumption_reward,
    COALESCE(operational_payments.total,0) as operational_payments,
    COALESCE(split_payments.total,0) as split_payments,
       COALESCE(payment.total, 0) + COALESCE(ai_payment.total, 0) + COALESCE(map_consumption_reward.total, 0) + COALESCE(operational_payments.total, 0) + COALESCE(split_payments.total, 0) as total from

    (select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
     from hivemapper.payments p
              inner join hivemapper.mints m on m.id = p.mint_id
              inner join hivemapper.transactions t on t.id = m.transaction_id
              inner join hivemapper.blocks b on b.id = t.block_id
     group by DATE_TRUNC('month', b.timestamp)) payment

        LEFT JOIN

    (select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
     from hivemapper.ai_payments p
              inner join hivemapper.mints m on m.id = p.mint_id
              inner join hivemapper.transactions t on t.id = m.transaction_id
              inner join hivemapper.blocks b on b.id = t.block_id
     group by DATE_TRUNC('month', b.timestamp)) ai_payment on payment.month = ai_payment.month

        LEFT JOIN

    (select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
     from hivemapper.map_consumption_reward p
              inner join hivemapper.mints m on m.id = p.mint_id
              inner join hivemapper.transactions t on t.id = m.transaction_id
              inner join hivemapper.blocks b on b.id = t.block_id
     group by DATE_TRUNC('month', b.timestamp)) map_consumption_reward on payment.month = map_consumption_reward.month

        LEFT JOIN

    (select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
     from hivemapper.operational_payments p
              inner join hivemapper.mints m on m.id = p.mint_id
              inner join hivemapper.transactions t on t.id = m.transaction_id
              inner join hivemapper.blocks b on b.id = t.block_id
     group by DATE_TRUNC('month', b.timestamp)) operational_payments on payment.month = operational_payments.month

        LEFT JOIN

    (select DATE_TRUNC('month', b.timestamp) as month, sum(m.amount) as total
     from hivemapper.split_payments p
              inner join hivemapper.mints m on m.id = p.fleet_mint_id or m.id = p.driver_mint_id
              inner join hivemapper.transactions t on t.id = m.transaction_id
              inner join hivemapper.blocks b on b.id = t.block_id
     group by DATE_TRUNC('month', b.timestamp)) split_payments on payment.month = split_payments.month

order by payment.month