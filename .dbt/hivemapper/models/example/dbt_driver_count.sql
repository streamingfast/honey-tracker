{{ config(materialized='table') }}

select count(*) as driver_count
from hivemapper.drivers