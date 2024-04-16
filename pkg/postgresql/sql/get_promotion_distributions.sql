-- Function: report.get_promotion_distributions
-- Description: Retrieves information about promotion distributions for a specific brand within a specified date range.

DROP FUNCTION IF EXISTS report.get_promotion_distributions(character varying, date, date);

CREATE OR REPLACE FUNCTION report.get_promotion_distributions(
    brand_code character varying,
    start_date date,
    end_date date
)
RETURNS TABLE(
    username character varying, 
    promotion_name character varying, 
    promotion_type text, 
    created_on timestamp with time zone, 
    bonus_amount numeric, 
    sent_on timestamp with time zone
)
LANGUAGE plpgsql
COST 100
VOLATILE
PARALLEL UNSAFE
ROWS 1000

AS $BODY$
BEGIN
    RETURN QUERY
    SELECT * FROM (
        SELECT 
            p.username,
            pm.name AS promotion_name,
            pm.promotion_type,
            pd.created_on,
            pd.bonus_amount,
            pd.sent_on
        FROM 
            dbo.promotion_distributes pd
            JOIN dbo.promotions pm ON pm.id = pd.promotion_id
            JOIN dbo.players p ON p.player_code = pd.player_code
        WHERE 
            p.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
            AND pd.sent_on >= start_date
            AND pd.sent_on < end_date

        UNION

        SELECT 
            p.username,
            pm.name AS promotion_name,
            pb.promotion_type,
            pb.created_on,
            pp.bonus_amount,
            pp.sent_on
        FROM 
            dbo.promotion_players pp 
            JOIN dbo.promotion_bonuses pb ON pb.id = pp.promotion_bonus_id
            JOIN dbo.promotions pm ON pm.id = pb.promotion_id
            JOIN dbo.players p ON p.player_code = pp.player_code
        WHERE 
            pb.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
            AND pp.sent_on >= start_date
            AND pp.sent_on < end_date
    ) a
    ORDER BY sent_on DESC;
END;
$BODY$;
