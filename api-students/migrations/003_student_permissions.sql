CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(20)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    name        VARCHAR(50)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name       VARCHAR(20) REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);

INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator dengan akses penuh'),
    ('staff', 'Staf pengelola data mahasiswa'),
    ('user', 'Pengguna biasa')
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (name, description) VALUES
    ('student:list', 'Melihat daftar seluruh mahasiswa'),
    ('student:read:any', 'Melihat detail data mahasiswa milik siapa saja'),
    ('student:create', 'Menambahkan data mahasiswa baru'),
    ('student:update:any', 'Mengubah data mahasiswa milik siapa saja'),
    ('student:delete', 'Menghapus data mahasiswa'),
    ('user:list', 'Melihat daftar seluruh user'),
    ('user:read:any', 'Melihat data user mana pun'),
    ('user:update:any', 'Mengubah data user mana pun'),
    ('user:delete', 'Menghapus user'),
    ('role:assign', 'Mengubah role milik user lain')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('admin', 'user:list'),
    ('admin', 'user:read:any'),
    ('admin', 'user:update:any'),
    ('admin', 'user:delete'),
    ('admin', 'role:assign'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create'),
    ('staff', 'user:list'),
    ('staff', 'user:read:any')
ON CONFLICT DO NOTHING;

UPDATE users SET role = 'user' WHERE role NOT IN (SELECT name FROM roles);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_role'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT fk_users_role
            FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE;
    END IF;
END $$;

ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM users LIMIT 1) THEN
        UPDATE students SET owner_id = (SELECT id FROM users ORDER BY id ASC LIMIT 1) WHERE owner_id IS NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_students_owner'
    ) THEN
        ALTER TABLE students ADD CONSTRAINT fk_students_owner
            FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT;
    END IF;
END $$;
