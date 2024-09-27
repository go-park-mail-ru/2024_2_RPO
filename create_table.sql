CREATE TABLE User (
    u_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE Board (
    b_id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES User(u_id)
);
CREATE TABLE User_to_Board (
    ub_id SERIAL PRIMARY KEY,
    u_id INT NOT NULL,
    b_id INT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at TIMESTAMP,
    added_by INT,
    updated_by INT,
    can_edit BOOLEAN DEFAULT FALSE,
    can_share BOOLEAN DEFAULT FALSE,
    can_invite_members BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    notification_level INT,
    FOREIGN KEY (u_id) REFERENCES User(u_id),
    FOREIGN KEY (b_id) REFERENCES Board(b_id),
    FOREIGN KEY (added_by) REFERENCES User(u_id),
    FOREIGN KEY (updated_by) REFERENCES User(u_id)
);
CREATE TABLE Column (
    col_id SERIAL PRIMARY KEY,
    b_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_index INT,
    FOREIGN KEY (b_id) REFERENCES Board(b_id)
);
CREATE TABLE Card (
    crd_id SERIAL PRIMARY KEY,
    col_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    order_index INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (col_id) REFERENCES Column(col_id)
);
CREATE TABLE CardUpdate (
    cu_id SERIAL PRIMARY KEY,
    crd_id INT NOT NULL,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    type VARCHAR(50),
    text TEXT,
    FOREIGN KEY (crd_id) REFERENCES Card(crd_id),
    FOREIGN KEY (created_by) REFERENCES User(u_id)
);
CREATE TABLE CardPinnedFile (
    cpf_id SERIAL PRIMARY KEY,
    crd_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    url VARCHAR(2048) NOT NULL,
    FOREIGN KEY (crd_id) REFERENCES Card(crd_id),
    FOREIGN KEY (created_by) REFERENCES User(u_id) }
);
CREATE TABLE Notification (
    u_id INT NOT NULL,
    b_id INT NOT NULL,
    type VARCHAR(50),
    text TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (u_id, b_id, created_at),
    FOREIGN KEY (u_id) REFERENCES User(u_id),
    FOREIGN KEY (b_id) REFERENCES Board(b_id)
);
