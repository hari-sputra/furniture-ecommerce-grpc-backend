INSERT INTO user_role (
    id, name, code, created_at, created_by, is_deleted
)

VALUES 
    (uuid_generate_v4(), 'admin', 1, NOW(), 'admin', FALSE),
    (uuid_generate_v4(), 'customer', 2, NOW(), 'admin', FALSE);