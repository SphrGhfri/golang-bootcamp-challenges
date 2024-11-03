-- Section1
SELECT 
    s.user_id, 
    ROUND(
        COALESCE(
            CASE WHEN COUNT(c.user_id) > 0 THEN 
                SUM(CASE WHEN c.action = 'confirmed' THEN 1 ELSE 0 END) * 1.0 / COUNT(c.user_id)
            ELSE 0 END, 
        0), 
        2
    ) AS confirmation_rate
FROM 
    signups s
LEFT JOIN 
    confirmations c ON s.user_id = c.user_id
GROUP BY 
    s.user_id
ORDER BY 
    s.user_id;
-- Section2
   SELECT 
    EXTRACT(HOUR FROM time_stamp) AS hour, 
    COUNT(*) AS timeout_count
FROM 
    confirmations
WHERE 
    action = 'timeout'  -- Filter for timeouts
GROUP BY 
    hour
ORDER BY 
    hour;
-- Section3
   SELECT 
    s.user_id,
    MAX(c.time_stamp) AS latest_confirmation
FROM 
    signups s
LEFT JOIN 
    confirmations c ON s.user_id = c.user_id
GROUP BY 
    s.user_id
ORDER BY 
    latest_confirmation ASC;