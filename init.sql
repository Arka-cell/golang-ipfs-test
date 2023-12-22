CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO messages (sender_id, receiver_id, text) VALUES
    (1, 2, 'Hello, how are you?'),
    (2, 1, 'I am good, thanks!'),
    (3, 1, 'Meeting at 3 PM today'),
    (1, 3, 'Sure, I will be there'),
    (2, 3, 'Don not forget to bring the documents'),
    (3, 2, 'Okay, got it');
