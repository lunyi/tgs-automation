-- Function: report.get_players_adjust_amount
-- Description: Retrieves adjustment transactions for a specified brand and date range, including details about the user and the transaction.

DROP FUNCTION IF EXISTS report.get_players_adjust_amount(character varying, date, date, integer);

CREATE OR REPLACE FUNCTION report.get_players_adjust_amount(
    brand_code character varying,
    p_start_date date,
    p_end_date date,
    p_type integer
)
RETURNS TABLE(
    "玩家用戶名" text,
    "公司調帳" numeric,
    "派發前餘額" numeric,
    "派發後餘額" numeric,
    "執行時間" timestamp with time zone,
    "執行者" text,
    "描述" text
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
        p.username::text AS "玩家用戶名", 
        t.amount AS "公司調帳", 
        t.balance - t.amount AS "派發前餘額", 
        t.balance AS "派發後餘額", 
        t.recorded_on AS "執行時間", 
        CASE
            WHEN u.id IS NULL THEN 'System'::text
            ELSE u.username::text
        END AS "執行者", 
        t.description::text AS "描述"
    FROM 
        dbo.transactions AS t
        INNER JOIN dbo.players AS p 
            ON t.player_code = p.player_code
        LEFT JOIN dbo.users AS u 
            ON t.created_by_user_id = u.id
    WHERE 
        t.brand_id = (SELECT id FROM dbo.brands WHERE code = brand_code)
        AND t.recorded_on >= p_start_date
        AND t.recorded_on < p_end_date
        AND t.type = p_type;
END;
$BODY$;
