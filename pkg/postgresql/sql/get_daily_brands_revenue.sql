CREATE OR REPLACE FUNCTION report.get_daily_brands_revenue(
    report_date_param date  -- 添加一个日期参数用于指定报告日期
)
RETURNS TABLE(platform text, currency_code text, date date, active_users_count numeric, daily_order_count text, daily_revenue_usd text, monthly_cumulative_revenue_usd text) 
LANGUAGE 'plpgsql'
COST 100
VOLATILE PARALLEL UNSAFE
ROWS 1000
AS $BODY$
DECLARE
    start_date date;
    end_date date;
BEGIN
    IF report_date_param = DATE_TRUNC('month', report_date_param) THEN
        -- 如果 report_date_param 是每个月的第一天，查询上个月的数据
        start_date := DATE_TRUNC('month', report_date_param - INTERVAL '1 month');
        end_date := DATE_TRUNC('month', report_date_param) - INTERVAL '1 day';
    ELSE
        -- 否则查询当前月的数据
        start_date := DATE_TRUNC('month', report_date_param);
        end_date := DATE_TRUNC('month', report_date_param) + INTERVAL '1 month' - INTERVAL '1 day';
    END IF;

    RETURN QUERY 
    WITH ExchangeRates AS (
            SELECT 'PHP' AS currency_code, 0.0177 AS rate_to_usd
            UNION ALL SELECT 'HKD', 0.1277
            UNION ALL SELECT 'VND_1000', 0.0398
            UNION ALL SELECT 'IDR_1000', 0.062
        ), SubQuery AS (
            SELECT 
                b.code AS platform,
                bc.currency_code,
                ((p.report_date::date + INTERVAL '8 hours')::timestamp)::date AS date,
                COUNT(DISTINCT p.player_code) AS active_users_count, 
                SUM(p.round_count) AS daily_order_count_raw,
                ROUND(SUM(p.win_loss_amount * er.rate_to_usd), 2) AS daily_revenue_usd_raw,
                SUM(ROUND(SUM(p.win_loss_amount * er.rate_to_usd),2)) OVER (
                    PARTITION BY b.code 
                    ORDER BY ((p.report_date::date + INTERVAL '8 hours')::timestamp)::date ASC
                ) AS monthly_cumulative_revenue_usd
            FROM report.player_aggregates p
            JOIN dbo.brands b ON b.id = p.brand_id
            JOIN dbo.brand_currencies bc ON b.id = bc.brand_id
            LEFT JOIN ExchangeRates er ON bc.currency_code = er.currency_code
            WHERE 
                 p.report_date >= start_date + INTERVAL '8 hours'
                AND p.report_date <= end_date + INTERVAL '1 day' + INTERVAL '8 hours'
                AND p.bet_amount > 0
                AND b.code NOT IN ('Sky8', 'GPI')
            GROUP BY b.code, bc.currency_code, date
        )

    SELECT 
        CAST('Total' AS text) AS platform,
        CAST('' AS text) AS currency_code,
        (report_date_param - INTERVAL '1 day')::date AS date,
        COALESCE(SUM(sq.active_users_count), 0) AS active_users_count,
        TO_CHAR(COALESCE(SUM(sq.daily_order_count_raw), 0), 'FM999,999,999,999') AS daily_order_count,
        TO_CHAR(COALESCE(SUM(sq.daily_revenue_usd_raw), 0), 'FM999,999,999,999.99') AS daily_revenue_usd,
        TO_CHAR(COALESCE(SUM(sq.monthly_cumulative_revenue_usd), 0), 'FM999,999,999,999.99') AS monthly_cumulative_revenue_usd
    FROM SubQuery sq
	WHERE sq.date = (report_date_param - INTERVAL '1 day')::date
    UNION ALL
    SELECT 
        CAST(sq.platform AS text),
        CAST(sq.currency_code AS text),
        (sq.date)::date AS date,
        COALESCE(sq.active_users_count, 0)  AS active_users_count,
        TO_CHAR(COALESCE(sq.daily_order_count_raw, 0), 'FM999,999,999,999') AS daily_order_count,
        TO_CHAR(COALESCE(sq.daily_revenue_usd_raw, 0), 'FM999,999,999,999.99') AS daily_revenue_usd,
        TO_CHAR(COALESCE(sq.monthly_cumulative_revenue_usd, 0), 'FM999,999,999,999.99') AS monthly_cumulative_revenue_usd
    FROM SubQuery sq
    WHERE sq.date = (report_date_param - INTERVAL '1 day')::date
    ORDER BY 2;
END; 
$BODY$;