package httpModel

type HttpForumHttpModel struct {
	Kind         string  `json:"kind"`
	Uuid         string  `json:"uuid"`
	Title        string  `json:"title"`
	UserId       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	UserColleage string  `json:"user_colleage"`
	UserIcon     string  `json:"user_icon"`
	CreatedTime  tString `json:"created_time"`
	ThumbUp      int     `json:"thumb_up"`
	Reply        int     `json:"reply"`
	Read         int     `json:"read"`
}
