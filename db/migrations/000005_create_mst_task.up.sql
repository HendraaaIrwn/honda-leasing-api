-- Schema: mst (master data)

-- 1. template_tasks <<mst>>
CREATE TABLE mst.template_tasks (
    teta_id      BIGSERIAL PRIMARY KEY,
    teta_name    VARCHAR(85) NOT NULL,
    teta_role_id BIGINT NOT NULL REFERENCES account.roles(role_id) ON DELETE RESTRICT,
    description  TEXT,
    sequence_no  SMALLINT DEFAULT 0,
    is_required  BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    call_function TEXT DEFAULT NULL
);

-- 2. template_task_attributes <<mst>>
CREATE TABLE mst.template_task_attributes (
    tetat_id       BIGSERIAL PRIMARY KEY,
    tetat_name     VARCHAR(85) NOT NULL,
    tetat_teta_id  BIGINT NOT NULL REFERENCES mst.template_tasks(teta_id) ON DELETE CASCADE,
    description    TEXT,
    is_required    BOOLEAN DEFAULT TRUE,
    attribute_type VARCHAR(50),
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index
CREATE INDEX idx_template_tasks_role ON mst.template_tasks(teta_role_id);
CREATE INDEX idx_template_tasks_seq  ON mst.template_tasks(sequence_no);
CREATE INDEX idx_tetat_teta          ON mst.template_task_attributes(tetat_teta_id);

-- Seed data
INSERT INTO mst.template_tasks (teta_name, teta_role_id, description, sequence_no, is_required) VALUES
('Input Pengajuan & Unggah Dokumen',
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'),
 'Isi form pengajuan, unggah KTP, KK, slip gaji/foto usaha',
 1, TRUE),
('Auto Scoring Awal & Pre-Approval',
 (SELECT role_id FROM account.roles WHERE role_name = 'ADMIN_CABANG'),
 'Sistem cek SLIK OJK, DSR, duplikasi data',
 2, TRUE),
('Survei Lapangan / Home Visit',
 (SELECT role_id FROM account.roles WHERE role_name = 'SURVEYOR'),
 'Kunjungi alamat, verifikasi, foto rumah, wawancara',
 3, TRUE),
('Input Hasil Survei & Rekomendasi',
 (SELECT role_id FROM account.roles WHERE role_name = 'SURVEYOR'),
 'Upload hasil survei & dokumen pendukung ke sistem',
 4, TRUE),
('Review & Approval Final (ACC/Reject)',
 (SELECT role_id FROM account.roles WHERE role_name = 'ADMIN_CABANG'),
 'Analisis hasil survei + scoring, beri keputusan ACC',
 5, TRUE),
('Akad & Tanda Tangan Kontrak',
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'),
 'Customer tanda tangan perjanjian, polis asuransi di dealer/cabang',
 6, TRUE),
('Pembayaran DP + Biaya Awal',
 (SELECT role_id FROM account.roles WHERE role_name = 'FINANCE'),
 'Verifikasi pembayaran DP, admin, asuransi, fidusia',
 7, TRUE),
('Proses PO & Pembelian Unit ke Dealer',
 (SELECT role_id FROM account.roles WHERE role_name = 'FINANCE'),
 'Leasing bayar ke dealer, terbitkan PO',
 8, TRUE),
('Delivery Motor ke Rumah Customer',
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'),
 'Unit dikirim, serah terima + dokumen (STNK sementara)',
 9, TRUE),
('Mulai Cicilan & Monitoring Pembayaran',
 (SELECT role_id FROM account.roles WHERE role_name = 'COLLECTION'),
 'Generate schedule angsuran, follow up pembayaran bulanan',
 10, TRUE),
('System Closed',
 (SELECT role_id FROM account.roles WHERE role_name = 'SYSTEM'),
 'System automatically closed',
 11, FALSE);

INSERT INTO mst.template_task_attributes (tetat_name, tetat_teta_id, description, is_required, attribute_type) VALUES
('Upload KTP', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), 'Foto KTP asli', TRUE, 'file'),
('Upload KK', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), 'Kartu Keluarga', TRUE, 'file'),
('Slip Gaji / Bukti Penghasilan', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), '3 bulan terakhir atau surat usaha', TRUE, 'file'),
('Foto Rumah Depan & Sekitar', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Minimal 3-5 foto', TRUE, 'file'),
('Foto Selfie + KTP di Lokasi', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Verifikasi identitas', TRUE, 'file'),
('Catatan Wawancara', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Pekerjaan, penghasilan, kondisi rumah', TRUE, 'text'),
('Tanda Tangan Digital / Fisik', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Akad & Tanda Tangan Kontrak'), 'Perjanjian pembiayaan', TRUE, 'file'),
('Foto Serah Terima Unit', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Delivery Motor ke Rumah Customer'), 'Customer terima motor + kunci', TRUE, 'file');
