-- Migration 004: Enforce unique email addresses for users

-- 1. Drop the existing non-unique email index
DROP INDEX IF EXISTS users_email_idx;

-- 2. Create a unique index on lowercase email (consistent with users_username_lower_idx)
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_idx ON users (LOWER(email));
