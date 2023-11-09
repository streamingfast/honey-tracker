{{ config(materialized='table') }}

select count(*) as fleet_count
from hivemapper.fleets