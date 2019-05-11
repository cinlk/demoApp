package httpModel

type HttpForumHttpModel struct {
	Kind         string  `json:"kind"`
	Uuid         string  `json:"uuid"`
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	UserId       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	UserColleage string  `json:"user_colleage"`
	UserIcon     string  `json:"user_icon"`
	CreatedTime  tString `json:"created_time"`
	ThumbUp      int     `json:"thumb_up"`
	Reply        int     `json:"reply"`
	Read         int     `json:"read"`
	// 自己是否收藏 和 点赞
	IsCollected bool `json:"is_collected"`
	IsLike      bool `json:"is_like"`
}

type HttpForumResponse struct {
	HttpResultModel
	Uuid string `json:"uuid"`
}

type HttpSubReplyInfo struct {
	ReplyId     string  `json:"reply_id"`
	UserIcon    string  `json:"user_icon"`
	UserName    string  `json:"user_name"`
	Colleage    string  `json:"colleage"`
	CreatedTime tString `json:"created_time"`
	Content     string  `json:"content"`
	LikeCount   int     `json:"like_count"`
	ReplyCount  int     `json:"reply_count"`
}
