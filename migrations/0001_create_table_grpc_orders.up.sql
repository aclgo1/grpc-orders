CREATE TYPE payment_method AS ENUM ('PIX', 'CREDIT_CARD', 'INTERNAL_BALANCE', 'BOLETO');
CREATE TYPE order_status AS ENUM ('PENDING', 'PAID', 'FAILED', 'CANCELLED', 'REFUNDED');
CREATE TYPE order_type AS ENUM ('PREMIUM_SUBSCRIPTION', 'BALANCE_DEPOSIT', 'PRODUCT_PURCHASE');

CREATE TABLE grpc_orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    type order_type NOT NULL,
    
    amount BIGINT NOT NULL, 
    
    payment_method payment_method NOT NULL,
    status order_status NOT NULL DEFAULT 'PENDING',
    
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb, 
    
    gateway_transaction_id VARCHAR(255),
    
    pix_qr_code TEXT,
    pix_expiration TIMESTAMP WITH TIME ZONE,
    
    card_token VARCHAR(255),
    card_expiration VARCHAR(7),
    
    boleto_url TEXT,
    boleto_barcode TEXT,
    boleto_expiration TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_grpc_orders_account ON grpc_orders(account_id);
CREATE INDEX idx_grpc_orders_status ON grpc_orders(status);