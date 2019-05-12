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
	UserId      string  `json:"user_id"`
	Colleage    string  `json:"colleage"`
	CreatedTime tString `json:"created_time"`
	Content     string  `json:"content"`
	LikeCount   int     `json:"like_count"`
	ReplyCount  int     `json:"reply_count"`
	IsLike      bool    `json:"is_like"`
}

// 二级回复
type HttpSecondReplyInfo struct {
	ReplyId        string  `json:"reply_id"`
	SecondReplyId  string  `json:"second_reply_id"`
	Content        string  `json:"content"`
	IsLike         bool    `json:"is_like"`
	CreatedTime    tString `json:"created_time"`
	LikeCount      int     `json:"like_count"`
	UserId         string  `json:"user_id"`
	UserIcon       string  `json:"user_icon"`
	UserName       string  `json:"user_name"`
	TalkedUserName string  `json:"talked_user_name"`
	TalkedUserId   string  `json:"talked_user_id"`
	ToHost         bool    `json:"to_host"`

	HostUserId string `json:"-"`
}
