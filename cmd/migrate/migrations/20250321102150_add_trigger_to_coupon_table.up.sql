-- Create the trigger function
CREATE OR REPLACE FUNCTION notify_new_coupon()
RETURNS TRIGGER AS $$
BEGIN
    -- Construct the payload with data from the new row
    PERFORM pg_notify(
        'new_coupon_channel',
        json_build_object(
            'id', NEW.id,
            'code', NEW.coupon_code,
            'discount', NEW.discount_percentage,
            'expiration_date', NEW.expiration_date,
            'is_active', NEW.is_active
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger on the coupons table
CREATE TRIGGER coupon_inserted_trigger
AFTER INSERT ON coupons
FOR EACH ROW
EXECUTE FUNCTION notify_new_coupon();