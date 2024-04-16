-- Function: report.get_registered_players
-- Description: Retrieves details about players registered within a specified date range for a given brand, including agent and host information.

DROP FUNCTION IF EXISTS report.get_registered_players_by_brand(character varying, timestamp with time zone, timestamp with time zone);

CREATE OR REPLACE FUNCTION report.get_registered_players_by_brand(
    brand_code character varying,
    start_date timestamp with time zone,
    end_date timestamp with time zone
)
RETURNS TABLE(
    agent text,
    host text,
    player text,
    real_name text,
    registered_on timestamp with time zone
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
        coalesce(a.username, '')::text AS agent,        -- Coalesce to handle possible NULL values
        coalesce(PIR.host, '')::text AS host,           -- Coalesce to handle possible NULL values
        p.username::text AS player,
        coalesce(p.real_name, '')::text AS real_name,   -- Coalesce to handle possible NULL values
        p.registered_on
    FROM 
        dbo.players p
        LEFT JOIN dbo.agents a ON a.id = p.agent_id
        LEFT JOIN dbo.player_ip_records AS PIR ON p.player_code = PIR.player_code
            AND PIR.ip_type = 1
    WHERE 
        p.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code) 
        AND p.registered_on >= start_date 
        AND p.registered_on < end_date
    ORDER BY 
        1 NULLS FIRST, 2, 3; -- Order by agent (NULLS FIRST), host, and player
END;
$BODY$;
