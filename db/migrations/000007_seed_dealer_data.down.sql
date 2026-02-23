DELETE FROM dealer.customers
WHERE nik IN ('3273010101010001', '3273010202020002');

DELETE FROM dealer.motor_assets
WHERE file_name IN (
    'beat_cbs_hitam_2026.jpg',
    'vario160_abs_putih_side.jpg',
    'stylo_krem_retro_front.jpg',
    'pcx160_hitam_3d_view.pdf',
    'adv160_hitam_action.jpg'
  );

DELETE FROM dealer.motors
WHERE nomor_rangka IN (
    'MH1JF41-001',
    'MH1JF41-002',
    'MH1JF41-003',
    'MH1JF41-004',
    'MH1JF50-001',
    'MH1JF50-002',
    'MH1JF60-001',
    'MH1JF60-002',
    'MH1JF70-001',
    'MH1JF70-002',
    'MH1JF80-001',
    'MH1JF90-001',
    'MH1JF100-001'
  );

DELETE FROM dealer.motor_types
WHERE moty_name IN ('Matic', 'Sport', 'Classic', 'Bebek', 'Maxi');
