CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    homeTeam VARCHAR(100) NOT NULL,
    awayTeam VARCHAR(100) NOT NULL,
    score1 INTEGER,
    score2 INTEGER,
    matchDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Datos de ejemplo
INSERT INTO matches (homeTeam, awayTeam, score1, score2) VALUES
('Real Madrid', 'Barcelona', 2, 6),
('Valencia', 'Villarreal', 2, 1),
('Real Sociedad', 'Athletic Club', 1, 1),
('Getafe', 'Betis', 0, 3),
('Espanyol', 'Celta de Vigo', 2, 2),
('Mallorca', 'Osasuna', 1, 0),
('Alavés', 'Granada', 3, 2),
('Levante', 'Rayo Vallecano', 1, 4),
('Cadiz', 'Elche', 0, 1),
('Sevilla', 'Barcelona', 2, 3),
('Real Madrid', 'Atletico Madrid', 2, 1);
