DELETE FROM account.role_permission rp
USING account.roles r
WHERE rp.role_id = r.role_id
  AND r.role_name IN ('SUPER_ADMIN', 'ADMIN_CABANG', 'SALES', 'SURVEYOR', 'COLLECTION', 'FINANCE', 'CUSTOMER', 'SYSTEM');

DELETE FROM account.user_roles ur
USING account.users u
WHERE ur.user_id = u.user_id
  AND u.phone_number IN (
    '+6281212345678',
    '+6281556789123',
    '+6281788990011',
    '+6282198765432',
    '+6285712345678',
    '+6289612345678',
    '+6281314151617',
    '+621'
  );

DELETE FROM account.user_oauth_provider uop
USING account.users u
WHERE uop.user_id = u.user_id
  AND u.phone_number IN (
    '+6281212345678',
    '+6281556789123',
    '+6281788990011',
    '+6282198765432',
    '+6285712345678',
    '+6289612345678',
    '+6281314151617',
    '+621'
  );

DELETE FROM account.oauth_providers
WHERE provider_name IN ('google', 'apple');

DELETE FROM account.users
WHERE phone_number IN (
    '+6281212345678',
    '+6281556789123',
    '+6281788990011',
    '+6282198765432',
    '+6285712345678',
    '+6289612345678',
    '+6281314151617',
    '+621'
  );

DELETE FROM account.permissions
WHERE permission_type IN (
    'view_dashboard',
    'view_contract',
    'create_contract',
    'approve_contract',
    'view_survey',
    'create_survey',
    'view_payment',
    'record_payment',
    'manage_user',
    'manage_oauth',
    'export_report',
    'send_notif'
  );

DELETE FROM account.roles
WHERE role_name IN ('SUPER_ADMIN', 'ADMIN_CABANG', 'SALES', 'SURVEYOR', 'COLLECTION', 'FINANCE', 'CUSTOMER', 'SYSTEM');
