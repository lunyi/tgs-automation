-- Function: report.get_first_deposited_players_by_brand
-- Description: Retrieves a list of players who made their first deposit within a specified date range for a given brand, including agent and IP host information.

DROP FUNCTION IF EXISTS report.get_first_deposited_players_by_brand(character varying, timestamp with time zone, timestamp with time zone);

CREATE OR REPLACE FUNCTION report.get_first_deposited_players_by_brand(
    brand_code character varying,
    start_date timestamp with time zone,
    end_date timestamp with time zone
)
RETURNS TABLE(
    agent text, 
    host text, 
    playername text, 
    daily_deposit_amount numeric, 
    daily_deposit_count integer, 
    first_deposit_on timestamp with time zone
)
LANGUAGE plpgsql
COST 100
VOLATILE
PARALLEL UNSAFE
ROWS 1000
AS $BODY$
BEGIN
    RETURN QUERY
    SELECT 
        coalesce(a.username, '')::text AS agent,
        coalesce(PIR.host, '')::text AS host,
        p.username::text AS playername, 
        pdp.daily_deposit_amount, 
        pdp.daily_deposit_count,
        pdp.first_deposit_on
    FROM 
        dbo.player_daily_payment_statistics pdp
        JOIN dbo.players p ON p.player_code = pdp.player_code
        LEFT JOIN dbo.player_ip_records PIR ON p.player_code = PIR.player_code
            AND PIR.ip_type = 1
        LEFT JOIN dbo.agents a ON a.id = p.agent_id
    WHERE 
        p.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
        AND pdp.first_deposit_on >= start_date 
        AND pdp.first_deposit_on < end_date 
    ORDER BY 
        agent NULLS FIRST,
        host,
        playername;
END;
$BODY$;
