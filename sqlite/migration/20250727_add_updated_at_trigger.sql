CREATE TRIGGER update_notes_updated_at
    AFTER UPDATE ON notes
    FOR EACH ROW
    WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
