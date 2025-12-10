-- Add trigger to prevent destructive changes to protected users
CREATE OR REPLACE FUNCTION prevent_protected_user_changes() RETURNS trigger AS $$
BEGIN
    IF OLD.is_protected THEN
        IF TG_OP = 'DELETE' THEN
            RAISE EXCEPTION 'Cannot delete a protected user';
        END IF;

        IF TG_OP = 'UPDATE' THEN
            -- Block role, email, protection flag change, and soft delete for protected users.
            IF NEW.is_protected IS DISTINCT FROM OLD.is_protected THEN
                RAISE EXCEPTION 'Cannot change is_protected for a protected user';
            END IF;
            IF NEW.role IS DISTINCT FROM OLD.role THEN
                RAISE EXCEPTION 'Cannot change role for a protected user';
            END IF;
            IF NEW.email IS DISTINCT FROM OLD.email THEN
                RAISE EXCEPTION 'Cannot change email for a protected user';
            END IF;
            IF NEW.deleted_at IS DISTINCT FROM OLD.deleted_at THEN
                RAISE EXCEPTION 'Cannot delete/soft-delete a protected user';
            END IF;
            -- Password updates and other benign fields remain allowed.
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
