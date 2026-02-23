DELETE FROM leasing.leasing_contract
WHERE contract_number IN ('KTR-2026-001', 'KTR-2026-002', 'KTR-2026-003');

DELETE FROM leasing.leasing_product
WHERE kode_produk IN ('DP-RINGAN-24', 'SUPER-KILAT-12', 'BEBAS-1TH-36', 'STANDAR-48', 'LOW-DP-36');
