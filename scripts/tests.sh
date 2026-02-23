#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080/leasing/api}"
CURL_BIN="${CURL_BIN:-curl}"
JQ_BIN="${JQ_BIN:-jq}"

command -v "$CURL_BIN" >/dev/null 2>&1 || {
  echo "Error: curl is required" >&2
  exit 1
}

command -v "$JQ_BIN" >/dev/null 2>&1 || {
  echo "Error: jq is required" >&2
  exit 1
}

RUN_TS="$(date +%s)"
RUN_RND="$(printf '%05d' $RANDOM)"
RUN_KEY="${RUN_TS}${RUN_RND}"

DATE_DOB="1995-05-05T00:00:00Z"
DATE_REQUEST="2026-01-10T00:00:00Z"
DATE_AKAD="2026-01-20T00:00:00Z"
DATE_CICIL="2026-02-20T00:00:00Z"
DATE_START="2026-01-11T00:00:00Z"
DATE_END="2026-01-25T00:00:00Z"
DATE_JATUH_TEMPO="2026-03-01T00:00:00Z"
DATE_BAYAR="2026-02-01T00:00:00Z"
DATE_DELIVERY="2026-02-15T00:00:00Z"

TRACE_FILE="$(mktemp)"
trap 'rm -f "$TRACE_FILE"' EXIT

LAST_BODY=""
LAST_STATUS=""

log() {
  printf '[INFO] %s\n' "$*" >&2
}

fail() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

mark_coverage() {
  local method="$1"
  local route_template="$2"
  printf '%s %s\n' "$method" "$route_template" >>"$TRACE_FILE"
}

api() {
  local method="$1"
  local path="$2"
  local expected_status="$3"
  local payload="${4:-}"
  local url="${BASE_URL}${path}"

  local response status body
  if [[ -n "$payload" ]]; then
    response="$($CURL_BIN -sS -X "$method" "$url" \
      -H 'Content-Type: application/json' \
      -d "$payload" \
      -w $'\n%{http_code}')"
  else
    response="$($CURL_BIN -sS -X "$method" "$url" -w $'\n%{http_code}')"
  fi

  status="${response##*$'\n'}"
  body="${response%$'\n'*}"

  LAST_STATUS="$status"
  LAST_BODY="$body"

  if [[ "$status" != "$expected_status" ]]; then
    printf '[ERROR] %s %s expected=%s got=%s\n' "$method" "$path" "$expected_status" "$status" >&2
    printf '[ERROR] Response body: %s\n' "$body" >&2
    exit 1
  fi

  if [[ "$expected_status" != "204" ]]; then
    local success
    success="$(printf '%s' "$body" | $JQ_BIN -r '.success // empty' 2>/dev/null || true)"
    if [[ "$success" != "true" ]]; then
      printf '[ERROR] %s %s returned success=%s\n' "$method" "$path" "$success" >&2
      printf '[ERROR] Response body: %s\n' "$body" >&2
      exit 1
    fi
  fi

  printf '[OK] %s %s [%s]\n' "$method" "$path" "$status" >&2
}

json_get() {
  local query="$1"
  printf '%s' "$LAST_BODY" | $JQ_BIN -r "$query"
}

require_value() {
  local value="$1"
  local label="$2"
  if [[ -z "$value" || "$value" == "null" ]]; then
    fail "Missing value for ${label}"
  fi
}

list_resource() {
  local endpoint="$1"
  mark_coverage "GET" "$endpoint"
  api "GET" "$endpoint" "200"
}

get_resource() {
  local endpoint="$1"
  local id="$2"
  mark_coverage "GET" "${endpoint}/:id"
  api "GET" "${endpoint}/${id}" "200"
}

create_resource() {
  local endpoint="$1"
  local payload="$2"
  mark_coverage "POST" "$endpoint"
  api "POST" "$endpoint" "201" "$payload"
}

update_resource() {
  local endpoint="$1"
  local id="$2"
  local payload="$3"
  mark_coverage "PUT" "${endpoint}/:id"
  api "PUT" "${endpoint}/${id}" "200" "$payload"
}

delete_resource() {
  local endpoint="$1"
  local id="$2"
  mark_coverage "DELETE" "${endpoint}/:id"
  api "DELETE" "${endpoint}/${id}" "200"
}

create_and_smoke_crud() {
  local endpoint="$1"
  local id_field="$2"
  local create_payload="$3"
  local update_payload="$4"

  create_resource "$endpoint" "$create_payload"
  local id
  id="$(json_get ".data.${id_field}")"
  require_value "$id" "${endpoint}.${id_field}"

  list_resource "$endpoint"
  get_resource "$endpoint" "$id"
  update_resource "$endpoint" "$id" "$update_payload"

  printf '%s\n' "$id"
}

workflow_post() {
  local endpoint="$1"
  local expected_status="$2"
  local payload="$3"

  mark_coverage "POST" "$endpoint"
  api "POST" "$endpoint" "$expected_status" "$payload"
}

extract_ids_by_field_match_int() {
  local endpoint="$1"
  local id_field="$2"
  local match_field="$3"
  local match_value="$4"

  list_resource "$endpoint"
  printf '%s' "$LAST_BODY" | $JQ_BIN -r \
    --arg idf "$id_field" \
    --arg mf "$match_field" \
    --argjson mv "$match_value" \
    '.data[]? | select(.[$mf] == $mv) | .[$idf]'
}

log "Testing endpoints on ${BASE_URL}"

# quick connectivity check
api "GET" "/account/roles" "200"

# -----------------------------
# ACCOUNT
# -----------------------------
ROLE_NAME="role_${RUN_KEY}"
ROLE_NAME_UPD="role_${RUN_KEY}_upd"
PERMISSION_TYPE="perm.${RUN_KEY}"
PERMISSION_TYPE_UPD="perm.${RUN_KEY}.upd"
USERNAME="user_${RUN_KEY}"
USERNAME_UPD="user_${RUN_KEY}_upd"
PHONE="08${RUN_KEY:0:10}"
EMAIL="${RUN_KEY}@example.com"
PROVIDER_NAME="provider_${RUN_KEY}"
PROVIDER_NAME_UPD="provider_${RUN_KEY}_upd"

ROLE_CREATE="$($JQ_BIN -nc --arg rn "$ROLE_NAME" '{role_name:$rn,description:"role for endpoint test"}')"
ROLE_UPDATE="$($JQ_BIN -nc --arg rn "$ROLE_NAME_UPD" '{role_name:$rn,description:"role updated"}')"
ROLE_ID="$(create_and_smoke_crud "/account/roles" "role_id" "$ROLE_CREATE" "$ROLE_UPDATE")"

PERMISSION_CREATE="$($JQ_BIN -nc --arg pt "$PERMISSION_TYPE" '{permission_type:$pt,description:"permission for endpoint test"}')"
PERMISSION_UPDATE="$($JQ_BIN -nc --arg pt "$PERMISSION_TYPE_UPD" '{permission_type:$pt,description:"permission updated"}')"
PERMISSION_ID="$(create_and_smoke_crud "/account/permissions" "permission_id" "$PERMISSION_CREATE" "$PERMISSION_UPDATE")"

ROLE_PERMISSION_CREATE="$($JQ_BIN -nc --argjson role_id "$ROLE_ID" --argjson permission_id "$PERMISSION_ID" '{role_id:$role_id,permission_id:$permission_id}')"
ROLE_PERMISSION_UPDATE="$($JQ_BIN -nc --argjson role_id "$ROLE_ID" --argjson permission_id "$PERMISSION_ID" '{role_id:$role_id,permission_id:$permission_id}')"
ROLE_PERMISSION_ID="$(create_and_smoke_crud "/account/role_permission" "role_permission_id" "$ROLE_PERMISSION_CREATE" "$ROLE_PERMISSION_UPDATE")"

USER_CREATE="$($JQ_BIN -nc \
  --arg username "$USERNAME" \
  --arg phone "$PHONE" \
  --arg email "$EMAIL" \
  '{username:$username,phone_number:$phone,email:$email,full_name:"Endpoint Tester",password:"Pass12345!",pin_key:"123456",is_active:true}')"
USER_UPDATE="$($JQ_BIN -nc --arg full_name "Endpoint Tester Updated" --arg password "Pass12345!updated" '{full_name:$full_name,password:$password}')"
USER_ID="$(create_and_smoke_crud "/account/users" "user_id" "$USER_CREATE" "$USER_UPDATE")"

OAUTH_CREATE="$($JQ_BIN -nc \
  --arg provider_name "$PROVIDER_NAME" \
  '{provider_name:$provider_name,client_id:"cid",client_secret:"secret",redirect_uri:"https://example.com/callback",issuer_url:"https://example.com",active:true}')"
OAUTH_UPDATE="$($JQ_BIN -nc \
  --arg provider_name "$PROVIDER_NAME_UPD" \
  '{provider_name:$provider_name,client_id:"cid-upd",client_secret:"secret-upd",redirect_uri:"https://example.com/callback-upd",issuer_url:"https://example.com/upd",active:true}')"
OAUTH_PROVIDER_ID="$(create_and_smoke_crud "/account/oauth_providers" "provider_id" "$OAUTH_CREATE" "$OAUTH_UPDATE")"

USER_OAUTH_CREATE="$($JQ_BIN -nc \
  --argjson user_id "$USER_ID" \
  --argjson provider_id "$OAUTH_PROVIDER_ID" \
  '{user_id:$user_id,provider_id:$provider_id,access_token:"access-token",refresh_token:"refresh-token",expires_at:"2027-01-01T00:00:00Z"}')"
USER_OAUTH_UPDATE="$($JQ_BIN -nc \
  --argjson user_id "$USER_ID" \
  --argjson provider_id "$OAUTH_PROVIDER_ID" \
  '{user_id:$user_id,provider_id:$provider_id,access_token:"access-token-upd",refresh_token:"refresh-token-upd"}')"
USER_OAUTH_ID="$(create_and_smoke_crud "/account/user_oauth_provider" "user_oauth_id" "$USER_OAUTH_CREATE" "$USER_OAUTH_UPDATE")"

USER_ROLE_CREATE="$($JQ_BIN -nc \
  --argjson user_id "$USER_ID" \
  --argjson role_id "$ROLE_ID" \
  --argjson assigned_by "$USER_ID" \
  '{user_id:$user_id,role_id:$role_id,assigned_by:$assigned_by}')"
USER_ROLE_UPDATE="$($JQ_BIN -nc \
  --argjson user_id "$USER_ID" \
  --argjson role_id "$ROLE_ID" \
  --argjson assigned_by "$USER_ID" \
  '{user_id:$user_id,role_id:$role_id,assigned_by:$assigned_by}')"
USER_ROLE_ID="$(create_and_smoke_crud "/account/user_roles" "user_role_id" "$USER_ROLE_CREATE" "$USER_ROLE_UPDATE")"

# -----------------------------
# MST
# -----------------------------
PROV_NAME="prov_${RUN_KEY}"
PROV_NAME_UPD="prov_${RUN_KEY}_upd"
KAB_NAME="kab_${RUN_KEY}"
KAB_NAME_UPD="kab_${RUN_KEY}_upd"
KEC_NAME="kec_${RUN_KEY}"
KEC_NAME_UPD="kec_${RUN_KEY}_upd"
KEL_NAME="kel_${RUN_KEY}"
KEL_NAME_UPD="kel_${RUN_KEY}_upd"

PROVINCE_CREATE="$($JQ_BIN -nc --arg prov_name "$PROV_NAME" '{prov_name:$prov_name}')"
PROVINCE_UPDATE="$($JQ_BIN -nc --arg prov_name "$PROV_NAME_UPD" '{prov_name:$prov_name}')"
PROV_ID="$(create_and_smoke_crud "/mst/province" "prov_id" "$PROVINCE_CREATE" "$PROVINCE_UPDATE")"

KABUPATEN_CREATE="$($JQ_BIN -nc --arg kab_name "$KAB_NAME" --argjson prov_id "$PROV_ID" '{kab_name:$kab_name,prov_id:$prov_id}')"
KABUPATEN_UPDATE="$($JQ_BIN -nc --arg kab_name "$KAB_NAME_UPD" --argjson prov_id "$PROV_ID" '{kab_name:$kab_name,prov_id:$prov_id}')"
KAB_ID="$(create_and_smoke_crud "/mst/kabupaten" "kab_id" "$KABUPATEN_CREATE" "$KABUPATEN_UPDATE")"

KECAMATAN_CREATE="$($JQ_BIN -nc --arg kec_name "$KEC_NAME" --argjson kab_id "$KAB_ID" '{kec_name:$kec_name,kab_id:$kab_id}')"
KECAMATAN_UPDATE="$($JQ_BIN -nc --arg kec_name "$KEC_NAME_UPD" --argjson kab_id "$KAB_ID" '{kec_name:$kec_name,kab_id:$kab_id}')"
KEC_ID="$(create_and_smoke_crud "/mst/kecamatan" "kec_id" "$KECAMATAN_CREATE" "$KECAMATAN_UPDATE")"

KELURAHAN_CREATE="$($JQ_BIN -nc --arg kel_name "$KEL_NAME" --argjson kec_id "$KEC_ID" '{kel_name:$kel_name,kec_id:$kec_id}')"
KELURAHAN_UPDATE="$($JQ_BIN -nc --arg kel_name "$KEL_NAME_UPD" --argjson kec_id "$KEC_ID" '{kel_name:$kel_name,kec_id:$kec_id}')"
KEL_ID="$(create_and_smoke_crud "/mst/kelurahan" "kel_id" "$KELURAHAN_CREATE" "$KELURAHAN_UPDATE")"

LOCATION_CREATE="$($JQ_BIN -nc --argjson kel_id "$KEL_ID" '{street_address:"Jl. Endpoint Test No.1",postal_code:"40211",longitude:"107.601234",latitude:"-6.901234",kel_id:$kel_id}')"
LOCATION_UPDATE="$($JQ_BIN -nc --argjson kel_id "$KEL_ID" '{street_address:"Jl. Endpoint Test No.2",postal_code:"40212",longitude:"107.601999",latitude:"-6.901999",kel_id:$kel_id}')"
LOCATION_ID="$(create_and_smoke_crud "/mst/locations" "location_id" "$LOCATION_CREATE" "$LOCATION_UPDATE")"

TEMPLATE_TASK_CREATE="$($JQ_BIN -nc --argjson role_id "$ROLE_ID" '{teta_name:"Scoring QA Task",teta_role_id:$role_id}')"
TEMPLATE_TASK_UPDATE="$($JQ_BIN -nc --argjson role_id "$ROLE_ID" '{teta_name:"Scoring QA Task Updated",teta_role_id:$role_id}')"
TETA_ID="$(create_and_smoke_crud "/mst/template_tasks" "teta_id" "$TEMPLATE_TASK_CREATE" "$TEMPLATE_TASK_UPDATE")"

TEMPLATE_TASK_ATTR_CREATE="$($JQ_BIN -nc --argjson teta_id "$TETA_ID" '{tetat_name:"attribute-qa",tetat_teta_id:$teta_id}')"
TEMPLATE_TASK_ATTR_UPDATE="$($JQ_BIN -nc --argjson teta_id "$TETA_ID" '{tetat_name:"attribute-qa-updated",tetat_teta_id:$teta_id}')"
TETAT_ID="$(create_and_smoke_crud "/mst/template_task_attributes" "tetat_id" "$TEMPLATE_TASK_ATTR_CREATE" "$TEMPLATE_TASK_ATTR_UPDATE")"

# -----------------------------
# DEALER
# -----------------------------
MOTY_NAME="moty_${RUN_KEY}"
MOTY_NAME_UPD="moty_${RUN_KEY}_upd"

MOTOR_TYPE_CREATE="$($JQ_BIN -nc --arg moty_name "$MOTY_NAME" '{moty_name:$moty_name}')"
MOTOR_TYPE_UPDATE="$($JQ_BIN -nc --arg moty_name "$MOTY_NAME_UPD" '{moty_name:$moty_name}')"
MOTY_ID="$(create_and_smoke_crud "/dealer/motor_types" "moty_id" "$MOTOR_TYPE_CREATE" "$MOTOR_TYPE_UPDATE")"

MOTOR_NORANGKA="NRG${RUN_KEY}A"
MOTOR_NOMESIN="NMS${RUN_KEY}A"
MOTOR_NOPOL="B${RUN_KEY:0:5}QA"

MOTOR_CREATE="$($JQ_BIN -nc \
  --argjson moty_id "$MOTY_ID" \
  --arg nomor_rangka "$MOTOR_NORANGKA" \
  --arg nomor_mesin "$MOTOR_NOMESIN" \
  --arg nomor_polisi "$MOTOR_NOPOL" \
  '{merk:"Honda",motor_type:"Matic",tahun:2025,warna:"Hitam",nomor_rangka:$nomor_rangka,nomor_mesin:$nomor_mesin,cc_mesin:"150",nomor_polisi:$nomor_polisi,status_unit:"ready",harga_otr:25000000,motor_moty_id:$moty_id}')"
MOTOR_UPDATE="$($JQ_BIN -nc '{warna:"Merah",status_unit:"ready"}')"
MOTOR_ID="$(create_and_smoke_crud "/dealer/motors" "motor_id" "$MOTOR_CREATE" "$MOTOR_UPDATE")"

MOTOR_ASSET_CREATE="$($JQ_BIN -nc --argjson motor_id "$MOTOR_ID" '{file_name:"unit.png",file_size:120.5,file_type:"png",file_url:"https://example.com/unit.png",moas_motor_id:$motor_id}')"
MOTOR_ASSET_UPDATE="$($JQ_BIN -nc --argjson motor_id "$MOTOR_ID" '{file_name:"unit-upd.png",file_size:121.5,file_type:"png",file_url:"https://example.com/unit-upd.png",moas_motor_id:$motor_id}')"
MOAS_ID="$(create_and_smoke_crud "/dealer/motor_assets" "moas_id" "$MOTOR_ASSET_CREATE" "$MOTOR_ASSET_UPDATE")"

CUSTOMER_NIK="1${RUN_KEY:0:15}"
CUSTOMER_PHONE="09${RUN_KEY:0:10}"
CUSTOMER_EMAIL="customer.${RUN_KEY}@example.com"

CUSTOMER_CREATE="$($JQ_BIN -nc \
  --arg nik "$CUSTOMER_NIK" \
  --arg phone "$CUSTOMER_PHONE" \
  --arg email "$CUSTOMER_EMAIL" \
  --arg dob "$DATE_DOB" \
  --argjson location_id "$LOCATION_ID" \
  '{nik:$nik,nama_lengkap:"Customer QA",tanggal_lahir:$dob,no_hp:$phone,email:$email,pekerjaan:"Karyawan",perusahaan:"PT QA",salary:7000000,location_id:$location_id}')"
CUSTOMER_UPDATE="$($JQ_BIN -nc --argjson location_id "$LOCATION_ID" '{nama_lengkap:"Customer QA Updated",pekerjaan:"Supervisor",salary:9000000,location_id:$location_id}')"
CUSTOMER_ID="$(create_and_smoke_crud "/dealer/customer" "customer_id" "$CUSTOMER_CREATE" "$CUSTOMER_UPDATE")"

# workflow dedicated motor + customer
WF_MOTOR_NORANGKA="NRG${RUN_KEY}B"
WF_MOTOR_NOMESIN="NMS${RUN_KEY}B"
WF_MOTOR_NOPOL="D${RUN_KEY:0:5}WF"
WF_MOTOR_CREATE="$($JQ_BIN -nc \
  --argjson moty_id "$MOTY_ID" \
  --arg nomor_rangka "$WF_MOTOR_NORANGKA" \
  --arg nomor_mesin "$WF_MOTOR_NOMESIN" \
  --arg nomor_polisi "$WF_MOTOR_NOPOL" \
  '{merk:"Honda",motor_type:"Matic",tahun:2026,warna:"Biru",nomor_rangka:$nomor_rangka,nomor_mesin:$nomor_mesin,cc_mesin:"160",nomor_polisi:$nomor_polisi,status_unit:"ready",harga_otr:30000000,motor_moty_id:$moty_id}')"
create_resource "/dealer/motors" "$WF_MOTOR_CREATE"
WF_MOTOR_ID="$(json_get '.data.motor_id')"
require_value "$WF_MOTOR_ID" "wf_motor_id"

WF_CUSTOMER_NIK="2${RUN_KEY:0:15}"
WF_CUSTOMER_PHONE="07${RUN_KEY:0:10}"
WF_CUSTOMER_EMAIL="wf.${RUN_KEY}@example.com"
WF_CUSTOMER_CREATE="$($JQ_BIN -nc \
  --arg nik "$WF_CUSTOMER_NIK" \
  --arg phone "$WF_CUSTOMER_PHONE" \
  --arg email "$WF_CUSTOMER_EMAIL" \
  --arg dob "$DATE_DOB" \
  --argjson location_id "$LOCATION_ID" \
  '{nik:$nik,nama_lengkap:"Customer WF",tanggal_lahir:$dob,no_hp:$phone,email:$email,pekerjaan:"Wiraswasta",perusahaan:"CV WF",salary:10000000,location_id:$location_id}')"
create_resource "/dealer/customer" "$WF_CUSTOMER_CREATE"
WF_CUSTOMER_ID="$(json_get '.data.customer_id')"
require_value "$WF_CUSTOMER_ID" "wf_customer_id"

# -----------------------------
# LEASING + PAYMENT
# -----------------------------
PRODUCT_CODE="KP${RUN_KEY:0:10}"
PRODUCT_CODE_UPD="KP${RUN_KEY:0:8}UP"

LEASING_PRODUCT_CREATE="$($JQ_BIN -nc --arg kode "$PRODUCT_CODE" '{kode_produk:$kode,nama_produk:"Produk QA",tenor_bulan:24,dp_persen_min:20,dp_persen_max:35,bunga_flat:6.5,admin_fee:350000,asuransi:true}')"
LEASING_PRODUCT_UPDATE="$($JQ_BIN -nc --arg kode "$PRODUCT_CODE_UPD" '{kode_produk:$kode,nama_produk:"Produk QA Updated",tenor_bulan:24,dp_persen_min:20,dp_persen_max:35,bunga_flat:6.8,admin_fee:360000,asuransi:true}')"
PRODUCT_ID="$(create_and_smoke_crud "/leasing/leasing_product" "product_id" "$LEASING_PRODUCT_CREATE" "$LEASING_PRODUCT_UPDATE")"

CONTRACT_NUM="CTR-${RUN_KEY}"
CONTRACT_NUM_UPD="CTR-${RUN_KEY}-UPD"
LEASING_CONTRACT_CREATE="$($JQ_BIN -nc \
  --arg contract_number "$CONTRACT_NUM" \
  --arg request_date "$DATE_REQUEST" \
  --arg tanggal_mulai_cicil "$DATE_CICIL" \
  --argjson customer_id "$CUSTOMER_ID" \
  --argjson motor_id "$MOTOR_ID" \
  --argjson product_id "$PRODUCT_ID" \
  '{contract_number:$contract_number,request_date:$request_date,tanggal_mulai_cicil:$tanggal_mulai_cicil,tenor_bulan:24,nilai_kendaraan:25000000,dp_dibayar:6000000,pokok_pinjaman:19000000,total_pinjaman:21000000,cicilan_per_bulan:875000,status:"draft",customer_id:$customer_id,motor_id:$motor_id,product_id:$product_id}')"
LEASING_CONTRACT_UPDATE="$($JQ_BIN -nc --arg contract_number "$CONTRACT_NUM_UPD" '{contract_number:$contract_number,status:"approved"}')"
CONTRACT_ID="$(create_and_smoke_crud "/leasing/leasing_contract" "contract_id" "$LEASING_CONTRACT_CREATE" "$LEASING_CONTRACT_UPDATE")"

LEASING_TASK_CREATE="$($JQ_BIN -nc \
  --arg startdate "$DATE_START" \
  --arg enddate "$DATE_END" \
  --argjson contract_id "$CONTRACT_ID" \
  --argjson role_id "$ROLE_ID" \
  '{task_name:"Task QA",startdate:$startdate,enddate:$enddate,sequence_no:1,status:"inprogress",contract_id:$contract_id,role_id:$role_id}')"
LEASING_TASK_UPDATE="$($JQ_BIN -nc '{task_name:"Task QA Updated",status:"completed",sequence_no:2}')"
TASK_ID="$(create_and_smoke_crud "/leasing/leasing_tasks" "task_id" "$LEASING_TASK_CREATE" "$LEASING_TASK_UPDATE")"

LEASING_TASK_ATTR_CREATE="$($JQ_BIN -nc --argjson task_id "$TASK_ID" '{tasa_name:"doc-check",tasa_value:"pending",tasa_status:"pending",tasa_leta_id:$task_id}')"
LEASING_TASK_ATTR_UPDATE="$($JQ_BIN -nc --argjson task_id "$TASK_ID" '{tasa_name:"doc-check-upd",tasa_value:"done",tasa_status:"completed",tasa_leta_id:$task_id}')"
TASA_ID="$(create_and_smoke_crud "/leasing/leasing_tasks_attributes" "tasa_id" "$LEASING_TASK_ATTR_CREATE" "$LEASING_TASK_ATTR_UPDATE")"

LEASING_DOC_CREATE="$($JQ_BIN -nc --argjson contract_id "$CONTRACT_ID" '{file_name:"kontrak.pdf",file_size:200.1,file_type:"pdf",file_url:"https://example.com/kontrak.pdf",contract_id:$contract_id}')"
LEASING_DOC_UPDATE="$($JQ_BIN -nc --argjson contract_id "$CONTRACT_ID" '{file_name:"kontrak-upd.pdf",file_size:201.1,file_type:"pdf",file_url:"https://example.com/kontrak-upd.pdf",contract_id:$contract_id}')"
LOC_ID="$(create_and_smoke_crud "/leasing/leasing_contract_documents" "loc_id" "$LEASING_DOC_CREATE" "$LEASING_DOC_UPDATE")"

PAYMENT_SCHEDULE_CREATE="$($JQ_BIN -nc --arg jatuh_tempo "$DATE_JATUH_TEMPO" --argjson contract_id "$CONTRACT_ID" '{angsuran_ke:1,jatuh_tempo:$jatuh_tempo,pokok:800000,margin:50000,total_tagihan:850000,status_pembayaran:"unpaid",contract_id:$contract_id}')"
PAYMENT_SCHEDULE_UPDATE="$($JQ_BIN -nc '{status_pembayaran:"paid"}')"
SCHEDULE_ID="$(create_and_smoke_crud "/payment/payment_schedule" "schedule_id" "$PAYMENT_SCHEDULE_CREATE" "$PAYMENT_SCHEDULE_UPDATE")"

PAYMENT_NUMBER="PAY-${RUN_KEY}"
PAYMENT_NUMBER_UPD="PAY-${RUN_KEY}-UPD"
PAYMENT_CREATE="$($JQ_BIN -nc \
  --arg nomor_bukti "$PAYMENT_NUMBER" \
  --arg tanggal_bayar "$DATE_BAYAR" \
  --argjson contract_id "$CONTRACT_ID" \
  --argjson schedule_id "$SCHEDULE_ID" \
  '{nomor_bukti:$nomor_bukti,jumlah_bayar:850000,tanggal_bayar:$tanggal_bayar,metode_pembayaran:"transfer",provider:"BCA",contract_id:$contract_id,schedule_id:$schedule_id}')"
PAYMENT_UPDATE="$($JQ_BIN -nc --arg nomor_bukti "$PAYMENT_NUMBER_UPD" '{nomor_bukti:$nomor_bukti,provider:"Mandiri"}')"
PAYMENT_ID="$(create_and_smoke_crud "/payment/payments" "payment_id" "$PAYMENT_CREATE" "$PAYMENT_UPDATE")"

# -----------------------------
# LEASING WORKFLOW
# -----------------------------
WF_SUBMIT_PAYLOAD="$($JQ_BIN -nc \
  --argjson customer_id "$WF_CUSTOMER_ID" \
  --argjson motor_id "$WF_MOTOR_ID" \
  --argjson product_id "$PRODUCT_ID" \
  --arg request_date "$DATE_REQUEST" \
  '{customer_id:$customer_id,motor_id:$motor_id,product_id:$product_id,dp_dibayar:7000000,tenor_bulan:24,request_date:$request_date,documents:[{file_name:"ktp.jpg",file_size:10.2,file_type:"jpg",file_url:"https://example.com/ktp.jpg"}] }')"
workflow_post "/leasing/workflow/submit-application" "201" "$WF_SUBMIT_PAYLOAD"
WF_CONTRACT_ID="$(json_get '.data.contract_id')"
require_value "$WF_CONTRACT_ID" "wf_contract_id"

WF_AUTOSCORING_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" '{contract_id:$contract_id,auto_approved:true,manual_review_ready:false,manual_approved:false,note:"auto-approved by qa script"}')"
workflow_post "/leasing/workflow/auto-scoring" "200" "$WF_AUTOSCORING_PAYLOAD"

WF_SURVEY_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" '{contract_id:$contract_id,decision:"approve",note:"survey passed"}')"
workflow_post "/leasing/workflow/survey" "200" "$WF_SURVEY_PAYLOAD"

WF_FINAL_APPROVAL_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" '{contract_id:$contract_id,approved:true,note:"final acc"}')"
workflow_post "/leasing/workflow/final-approval" "200" "$WF_FINAL_APPROVAL_PAYLOAD"

WF_AKAD_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" --arg akad_date "$DATE_AKAD" --arg tanggal_mulai_cicil "$DATE_CICIL" '{contract_id:$contract_id,akad_date:$akad_date,tanggal_mulai_cicil:$tanggal_mulai_cicil,generate_contract_code:true}')"
workflow_post "/leasing/workflow/akad" "200" "$WF_AKAD_PAYLOAD"

WF_PAYMENT_NUMBER="PAY-WF-${RUN_KEY}"
WF_INITIAL_PAYMENT_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" --arg nomor_bukti "$WF_PAYMENT_NUMBER" --arg tanggal_bayar "$DATE_BAYAR" '{contract_id:$contract_id,nomor_bukti:$nomor_bukti,jumlah_bayar:7000000,tanggal_bayar:$tanggal_bayar,metode_pembayaran:"transfer",provider:"BCA"}')"
workflow_post "/leasing/workflow/initial-payment" "200" "$WF_INITIAL_PAYMENT_PAYLOAD"

WF_DEALER_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" '{contract_id:$contract_id,unit_ready_stock:true,estimated_indent_week:0,note:"stock ready"}')"
workflow_post "/leasing/workflow/dealer-fulfillment" "200" "$WF_DEALER_PAYLOAD"

WF_DELIVERY_PAYLOAD="$($JQ_BIN -nc --argjson contract_id "$WF_CONTRACT_ID" --arg delivery_date "$DATE_DELIVERY" '{contract_id:$contract_id,delivery_date:$delivery_date,customer_received:true,document_handover:true,handover_note:"dokumen diserahkan",contract_doc_uploads:[{file_name:"bast.pdf",file_size:11.2,file_type:"pdf",file_url:"https://example.com/bast.pdf"}]}')"
workflow_post "/leasing/workflow/delivery" "200" "$WF_DELIVERY_PAYLOAD"

# ensure workflow contract detail endpoint still works
get_resource "/leasing/leasing_contract" "$WF_CONTRACT_ID"

# -----------------------------
# DELETE / CLEANUP (reverse dependencies)
# -----------------------------

# payment related

delete_resource "/payment/payments" "$PAYMENT_ID"
for id in $(extract_ids_by_field_match_int "/payment/payments" "payment_id" "contract_id" "$WF_CONTRACT_ID"); do
  delete_resource "/payment/payments" "$id"
done

for id in $(extract_ids_by_field_match_int "/payment/payment_schedule" "schedule_id" "contract_id" "$WF_CONTRACT_ID"); do
  delete_resource "/payment/payment_schedule" "$id"
done

delete_resource "/payment/payment_schedule" "$SCHEDULE_ID"

# leasing task attrs

delete_resource "/leasing/leasing_tasks_attributes" "$TASA_ID"

for task_id in $(extract_ids_by_field_match_int "/leasing/leasing_tasks" "task_id" "contract_id" "$WF_CONTRACT_ID"); do
  for attr_id in $(extract_ids_by_field_match_int "/leasing/leasing_tasks_attributes" "tasa_id" "tasa_leta_id" "$task_id"); do
    delete_resource "/leasing/leasing_tasks_attributes" "$attr_id"
  done
done

# leasing tasks

delete_resource "/leasing/leasing_tasks" "$TASK_ID"
for id in $(extract_ids_by_field_match_int "/leasing/leasing_tasks" "task_id" "contract_id" "$WF_CONTRACT_ID"); do
  delete_resource "/leasing/leasing_tasks" "$id"
done

# leasing documents

delete_resource "/leasing/leasing_contract_documents" "$LOC_ID"
for id in $(extract_ids_by_field_match_int "/leasing/leasing_contract_documents" "loc_id" "contract_id" "$WF_CONTRACT_ID"); do
  delete_resource "/leasing/leasing_contract_documents" "$id"
done

# leasing contracts

delete_resource "/leasing/leasing_contract" "$CONTRACT_ID"
delete_resource "/leasing/leasing_contract" "$WF_CONTRACT_ID"

# product

delete_resource "/leasing/leasing_product" "$PRODUCT_ID"

# dealer
delete_resource "/dealer/motor_assets" "$MOAS_ID"
delete_resource "/dealer/motors" "$MOTOR_ID"
delete_resource "/dealer/motors" "$WF_MOTOR_ID"
delete_resource "/dealer/customer" "$CUSTOMER_ID"
delete_resource "/dealer/customer" "$WF_CUSTOMER_ID"
delete_resource "/dealer/motor_types" "$MOTY_ID"

# mst template + wilayah
delete_resource "/mst/template_task_attributes" "$TETAT_ID"
delete_resource "/mst/template_tasks" "$TETA_ID"
delete_resource "/mst/locations" "$LOCATION_ID"
delete_resource "/mst/kelurahan" "$KEL_ID"
delete_resource "/mst/kecamatan" "$KEC_ID"
delete_resource "/mst/kabupaten" "$KAB_ID"
delete_resource "/mst/province" "$PROV_ID"

# account
delete_resource "/account/user_roles" "$USER_ROLE_ID"
delete_resource "/account/user_oauth_provider" "$USER_OAUTH_ID"
delete_resource "/account/role_permission" "$ROLE_PERMISSION_ID"
delete_resource "/account/oauth_providers" "$OAUTH_PROVIDER_ID"
delete_resource "/account/users" "$USER_ID"
delete_resource "/account/permissions" "$PERMISSION_ID"
delete_resource "/account/roles" "$ROLE_ID"

# -----------------------------
# COVERAGE CHECK
# -----------------------------

CRUD_RESOURCES=(
  "/account/oauth_providers"
  "/account/users"
  "/account/user_oauth_provider"
  "/account/roles"
  "/account/user_roles"
  "/account/permissions"
  "/account/role_permission"
  "/mst/province"
  "/mst/kabupaten"
  "/mst/kecamatan"
  "/mst/kelurahan"
  "/mst/locations"
  "/mst/template_tasks"
  "/mst/template_task_attributes"
  "/dealer/motor_types"
  "/dealer/motors"
  "/dealer/motor_assets"
  "/dealer/customer"
  "/leasing/leasing_product"
  "/leasing/leasing_contract"
  "/leasing/leasing_tasks"
  "/leasing/leasing_tasks_attributes"
  "/leasing/leasing_contract_documents"
  "/payment/payment_schedule"
  "/payment/payments"
)

WORKFLOW_ENDPOINTS=(
  "/leasing/workflow/submit-application"
  "/leasing/workflow/auto-scoring"
  "/leasing/workflow/survey"
  "/leasing/workflow/final-approval"
  "/leasing/workflow/akad"
  "/leasing/workflow/initial-payment"
  "/leasing/workflow/dealer-fulfillment"
  "/leasing/workflow/delivery"
)

for resource in "${CRUD_RESOURCES[@]}"; do
  for requirement in \
    "GET ${resource}" \
    "GET ${resource}/:id" \
    "POST ${resource}" \
    "PUT ${resource}/:id" \
    "DELETE ${resource}/:id"; do
    if ! grep -Fxq "$requirement" "$TRACE_FILE"; then
      fail "Coverage missing: ${requirement}"
    fi
  done
done

for endpoint in "${WORKFLOW_ENDPOINTS[@]}"; do
  requirement="POST ${endpoint}"
  if ! grep -Fxq "$requirement" "$TRACE_FILE"; then
    fail "Coverage missing: ${requirement}"
  fi
done

log "All endpoint tests passed and coverage is complete."
