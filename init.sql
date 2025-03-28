CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    team1 VARCHAR(100) NOT NULL,
    team2 VARCHAR(100) NOT NULL,
    score1 INTEGER,
    score2 INTEGER,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Datos de ejemplo
INSERT INTO matches (team1, team2, score1, score2) VALUES
('Barcelona', 'Real Madrid', 2, 2),
('Atletico Madrid', 'Sevilla', 1, 0);