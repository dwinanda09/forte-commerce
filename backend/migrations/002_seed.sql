-- +goose Up
-- Password: admin123 | bcrypt cost 12 | per-user salt
INSERT INTO users (username, password_hash, password_salt)
VALUES (
    'admin',
    '$2a$12$6h1.RXCXtvmP.cjBdQfBZOJvKBxe2dZPwxIKJHvFdnKWJGT6cfn2W',
    '0000000000000000000000000000000000000000000000000000000000000000'
) ON CONFLICT (username) DO NOTHING;

-- Demo users: username = password (demo1..demo4) | bcrypt cost 12 | zero salt
INSERT INTO users (username, password_hash, password_salt)
VALUES
    ('demo1', '$2a$12$AZm.sJU2rSvSfxjs.QiS0uWpQvmd1A.03sRYp/h0J9OWB13HzteuS', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo2', '$2a$12$Fb8ipre94YraFBlne51e7OpJxedAwGhmHJ7q4MFwH9d7izAeVH1fO', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo3', '$2a$12$VIglvK2NcFJyRKyDAOTC2OcYboLZVCDHRqq/HreHSHGoHb0.iHglC', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo4', '$2a$12$idMlFvRkCNiT1ZIOPQjMO.g6Vt7.foBl3nLR9gHRjqUuHkkG380zy', '0000000000000000000000000000000000000000000000000000000000000000')
ON CONFLICT (username) DO NOTHING;

INSERT INTO products (sku, name, price, inventory_qty)
VALUES
    ('120P90', 'Google Home',    49.99,   10),
    ('43N23P', 'MacBook Pro',    5399.99,  5),
    ('A304SD', 'Alexa Speaker',  109.50,  10),
    ('234234', 'Raspberry Pi B',  30.00,   2)
ON CONFLICT (sku) DO NOTHING;

-- +goose Down
DELETE FROM products WHERE sku IN ('120P90', '43N23P', 'A304SD', '234234');
DELETE FROM users WHERE username IN ('admin', 'demo1', 'demo2', 'demo3', 'demo4');
