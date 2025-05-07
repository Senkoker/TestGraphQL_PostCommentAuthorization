package PostgresHandler

const (
	getUserInfo = "SELECT user_id,first_name,second_name,img_url FROM users_info WHERE user_id = ANY ($1)"
)
