DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS sharing_info;
CREATE TABLE users (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash BLOB NOT NULL
);
CREATE TABLE notes (
    id TEXT PRIMARY KEY NOT NULL,
    author_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_notes_author_id ON notes(author_id);
CREATE TABLE sharing_info (
    id TEXT PRIMARY KEY NOT NULL,
    note_id TEXT NOT NULL,
    editable BOOLEAN NOT NULL,
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);

CREATE TRIGGER update_notes_updated_at
    AFTER UPDATE ON notes
    FOR EACH ROW
    WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
