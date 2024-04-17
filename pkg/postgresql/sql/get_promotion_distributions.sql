-- Function: report.get_promotion_distributions
-- Description: Retrieves information about promotion distributions for a specific brand within a specified date range.

DROP FUNCTION IF EXISTS report.get_promotion_distributions(character varying, timestamp with time zone, timestamp with time zone);

CREATE OR REPLACE FUNCTION report.get_promotion_distributions(
	brand_code character varying,
	start_date timestamp with time zone,
	end_date timestamp with time zone)
    RETURNS 
TABLE(
	username character varying, 
	promotion_name character varying, 
	promotion_type text, 
	bonus_amount numeric, 
	created_on text, 
	sent_on text) 
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
BEGIN
    RETURN QUERY
    SELECT 
		a.username,
        a.name,
        a.promotion_type,
        a.bonus_amount,
	    to_char(a.created_on + interval '8 hours', 'HH24:MI:SS DD/MM/YYYY') AS created_on,
		COALESCE(to_char(a.sent_on + interval '8 hours', 'HH24:MI:SS DD/MM/YYYY'), '') AS sent_on
	FROM (
        SELECT 
            p.username,
            pm.name,
            pm.promotion_type,
            pd.bonus_amount,
		    pd.created_on,
		    pd.sent_on
        FROM dbo.promotion_distributes pd
        JOIN dbo.promotions pm ON pm.id = pd.promotion_id
        JOIN dbo.players p ON p.player_code = pd.player_code
        WHERE p.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
            AND pd.sent_on >= start_date
            AND pd.sent_on <  end_date
        UNION
        SELECT 
            p.username,
            pm.name,
            pb.promotion_type,
            pp.bonus_amount,
		    pb.created_on,
		    pp.sent_on
        FROM dbo.promotion_players pp 
        JOIN dbo.promotion_bonuses pb ON pb.id = pp.promotion_bonus_id
        JOIN dbo.promotions pm ON pm.id = pb.promotion_id
        JOIN dbo.players p ON p.player_code = pp.player_code
        WHERE pb.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
            AND pp.sent_on >=  start_date
            AND pp.sent_on < end_date
    ) a
    ORDER BY a.sent_on DESC;
END;
$BODY$;