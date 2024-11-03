-- Section1
   SELECT 
    c.name AS club_name,
    COUNT(p.player_id) AS foreign_player_count
FROM 
    players p
JOIN 
    clubs c ON p.current_club_id = c.club_id
JOIN 
    competitions comp ON c.domestic_competition_id = comp.competition_id
WHERE 
    p.country_of_citizenship != comp.country_name
GROUP BY 
    c.name
ORDER BY 
    foreign_player_count DESC, 
    club_name ASC;
-- Section2
WITH CompetitionEvents AS (
    SELECT 
        comp.name AS name,
        COUNT(ge.game_event_id) AS events
    FROM 
        competitions comp
    LEFT JOIN 
        games g ON comp.competition_id = g.competition_id
    LEFT JOIN 
        game_events ge ON g.game_id = ge.game_id
    GROUP BY 
        comp.competition_id, comp.name
    HAVING 
        COUNT(ge.game_event_id) > 0
),
RankedCompetitions AS (
    SELECT 
        name,
        events,
        DENSE_RANK() OVER (ORDER BY events DESC) AS ranking
    FROM 
        CompetitionEvents
)
SELECT 
    ranking,
    name,
    events
FROM 
    RankedCompetitions
ORDER BY 
    ranking ASC, name ASC;
-- Section3
SELECT 
    p.player_code,
    p.contract_expiration_date,
    SUM(a.goals) AS total_goals,
    SUM(a.assists) AS total_assists,
    SUM(a.minutes_played) AS total_minutes
FROM 
    appearances a
JOIN 
    players p ON a.player_id = p.player_id
JOIN 
    games g ON a.game_id = g.game_id
WHERE 
    g.season = (SELECT MAX(season) FROM games)
    AND p.contract_expiration_date < '2025-01-01'
GROUP BY 
    p.player_code, p.contract_expiration_date
ORDER BY 
    p.contract_expiration_date ASC, total_goals DESC, total_assists DESC, total_minutes DESC;
