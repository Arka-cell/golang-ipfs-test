package Tables

// Message struct represents the structure of the 'messages' table
type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	IPNS       string `json:"ipns"`
	IPFS       string `json:"ipfs"`
}
