-- Function: report.report_get_daily_brands_revenue
-- Description: This function retrieves daily revenue aggregates for brands, 
--              adjusting for different currencies and consolidating data by platform.

--DROP FUNCTION IF EXISTS report.get_daily_brands_revenue();

CREATE OR REPLACE FUNCTION report.get_daily_brands_revenue()
RETURNS TABLE(
    platform text, 
    currency_code text, 
    date date, 
    active_users_count numeric, 
    daily_order_count text, 
    daily_revenue_usd text, 
    monthly_cumulative_revenue_usd text
)
LANGUAGE plpgsql
COST 100
VOLATILE
PARALLEL UNSAFE
ROWS 1000
AS $BODY$
BEGIN
    RETURN QUERY 
    WITH ExchangeRates AS (
        SELECT 'PHP' AS currency_code, 0.0177 AS rate_to_usd
        UNION ALL SELECT 'HKD', 0.1277
        UNION ALL SELECT 'VND_1000', 0.0398
    ), SubQuery AS (
        SELECT 
            b.code AS platform,
            bc.currency_code,
            ((p.report_date::date + INTERVAL '1 day' + INTERVAL '8 hours')::timestamp)::date AS date,
            COUNT(DISTINCT p.player_code) AS active_users_count, 
            SUM(p.round_count) AS daily_order_count_raw,
            ROUND(SUM(p.win_loss_amount * er.rate_to_usd), 2) AS daily_revenue_usd_raw,
            SUM(ROUND(SUM(p.win_loss_amount * er.rate_to_usd), 2)) OVER (
                PARTITION BY b.code 
                ORDER BY ((p.report_date::date + INTERVAL '1 day' + INTERVAL '8 hours')::timestamp)::date ASC
            ) AS monthly_cumulative_revenue_usd
        FROM 
            report.player_aggregates p
            JOIN dbo.brands b ON b.id = p.brand_id
            JOIN dbo.brand_currencies bc ON b.id = bc.brand_id
            LEFT JOIN ExchangeRates er ON bc.currency_code = er.currency_code
        WHERE 
            p.report_date >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 day' + INTERVAL '8 hours'
            AND p.report_date < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month' + INTERVAL '8 hours'
            AND p.bet_amount > 0
            AND b.code NOT IN ('Sky8', 'GPI')
        GROUP BY 
            b.code, bc.currency_code, date
    )
    SELECT 
        CAST('Total' AS text) AS platform,
        CAST('' AS text) AS currency_code,
        (CURRENT_DATE - INTERVAL '1 day')::date AS date,
	    COALESCE(SUM(sq.active_users_count), 0) AS active_users_count,
		TO_CHAR(COALESCE(SUM(sq.daily_order_count_raw), 0), 'FM999,999,999,999') AS daily_order_count,
		TO_CHAR(COALESCE(SUM(sq.daily_revenue_usd_raw), 0), 'FM999,999,999,999.99') AS daily_revenue_usd,
		TO_CHAR(COALESCE(SUM(sq.monthly_cumulative_revenue_usd), 0), 'FM999,999,999,999.99') AS monthly_cumulative_revenue_usd
    FROM 
        SubQuery sq
    --WHERE 
      --  sq.date = (CURRENT_DATE - INTERVAL '1 day')::date
    UNION ALL
    SELECT 
        sq.platform,
        sq.currency_code,
		(sq.date - INTERVAL '1 day')::date AS date,
		COALESCE(sq.active_users_count, 0)  AS active_users_count,
		TO_CHAR(COALESCE(sq.daily_order_count_raw, 0), 'FM999,999,999,999') AS daily_order_count,
		TO_CHAR(COALESCE(sq.daily_revenue_usd_raw, 0), 'FM999,999,999,999.99') AS daily_revenue_usd,
		TO_CHAR(COALESCE(sq.monthly_cumulative_revenue_usd, 0), 'FM999,999,999,999.99') AS monthly_cumulative_revenue_usd
    FROM 
        SubQuery sq
    --WHERE 
      --  sq.date = (CURRENT_DATE - INTERVAL '1 day')::date
    ORDER BY 
        2;
END; 
$BODY$;
