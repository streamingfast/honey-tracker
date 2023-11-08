
/*
    Welcome to your first dbt model!
    Did you know that you can also configure models directly within SQL files?
    This will override configurations stated in dbt_project.yml

    Try changing "table" to "view" below
*/

{{ config(materialized='view') }}

select sum(mints.amount) as total_payments from hivemapper.payments
inner join hivemapper.mints on mints.id = payments.mint_id

/*
    Uncomment the line below to remove records with null `id` values
*/

-- where id is not null
