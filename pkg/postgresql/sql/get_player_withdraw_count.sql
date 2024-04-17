

--DROP FUNCTION IF EXISTS report.get_player_withdraw_count(character varying, date, date);
--select * from report.get_player_withdraw_count('MOPH', '20240401', '20240410')

CREATE OR REPLACE FUNCTION report.get_player_withdraw_count(
    brand_code VARCHAR,
    start_date DATE,
    end_date DATE
)
RETURNS INTEGER AS $$
DECLARE
    withdraw_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO withdraw_count
    FROM (
        SELECT pdps.player_code, COUNT(*) AS cc
        FROM dbo.player_daily_payment_statistics pdps
        JOIN dbo.players p ON p.player_code = pdps.player_code
        JOIN dbo.brands b ON p.brand_id = b.id
        WHERE b.code = brand_code
          AND pdps.report_date >= start_date
          AND pdps.report_date < end_date
          AND pdps.daily_withdraw_amount > 0
        GROUP BY pdps.player_code
    ) a
    WHERE a.cc > 0;
	RETURN withdraw_count;
END;
$$ LANGUAGE plpgsql;