CREATE TABLE IF NOT EXISTS students ( 
    id           SERIAL       PRIMARY KEY, 
    nim          CHAR(8)      NOT NULL UNIQUE, 
    name         VARCHAR(100) NOT NULL, 
    grade        NUMERIC(3,2) CHECK (grade >= 0 AND grade <= 4.00), 
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE, 
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW() 
);