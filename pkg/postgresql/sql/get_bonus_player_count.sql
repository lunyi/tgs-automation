-- FUNCTION: report.get_bonus_player_count(character varying, timestamp with time zone, timestamp with time zone)

-- DROP FUNCTION IF EXISTS report.get_bonus_player_count(character varying, timestamp with time zone, timestamp with time zone);

-- SELECT report.get_bonus_player_count('MOPH', '20240401+8', '20240410+8');

CREATE OR REPLACE FUNCTION report.get_bonus_player_count(
    brand_code VARCHAR,
    start_date timestamp with time zone,
    end_date timestamp with time zone
)
RETURNS INTEGER AS $$
DECLARE
    player_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO player_count
    FROM (
        SELECT 
            player_code,
            SUM(birthday_amount + season_amount + upgrade_amount + promotion_amount + rebate_amount) AS bonus
        FROM 
            report.player_aggregates 
        WHERE 
            brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
            AND report_date >= start_date
            AND report_date < end_date 
        GROUP BY 
            player_code
    ) a 
    WHERE a.bonus > 0;

    RETURN player_count;
END;
$$ LANGUAGE plpgsql;