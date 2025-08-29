CREATE TABLE IF NOT EXISTS tb_r_shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tracking_number TEXT NOT NULL UNIQUE,            -- nomor resi
    user_id UUID NOT NULL,                           -- ID user (dari auth service, tidak FK lintas DB)
    sender_name TEXT NOT NULL,
    sender_address TEXT NOT NULL,
    receiver_name TEXT NOT NULL,
    receiver_address TEXT NOT NULL,
    item_description TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('Created','Shipped','InTransit','Delivered','Cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_shipments_user ON tb_r_shipments(user_id);
CREATE INDEX IF NOT EXISTS idx_shipments_tracking ON tb_r_shipments(tracking_number);