{{ config(materialized='table') }}

select (SELECT SUM(amount) FROM hivemapper.mints) - (SELECT COALESCE(SUM(amount), 0)  FROM hivemapper.burns) as total_supply