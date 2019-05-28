package main

import (
	"demoApp/server/model/dbModel"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"strconv"
	"time"
)

var testDB *gorm.DB

type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type gistModel struct {
	gorm.Model
	Name  string
	Point point
}

func loadDB() {

	db, err := gorm.Open("postgres",
		"host=localhost user=wxadmin dbname=app sslmode=disable password=wxadmin")

	if err != nil {
		fmt.Println(err)
	}

	testDB = db
	testDB.LogMode(true)
	//AddDepartment()
	testDB.SingularTable(true)
	testDB.Set("gorm:table_options", "charset=utf8")
	testDB.DB().SetMaxIdleConns(3)
	testDB.DB().SetMaxIdleConns(1)

	err = testDB.AutoMigrate(gistModel{}).Error
	if err != nil {
		fmt.Print(err)
	}

}

func createInternJobs() {
	for i := 0; i < 10; i++ {
		item := dbModel.InternJobs{}
		item.Id = "intern-jobs-" + strconv.Itoa(i)
		item.Type = "intern"
		item.Name = "研发订单"
		item.CompanyID = "companyID"
		item.IconURL = "http://icons.iconarchive.com/icons/blackvariant/button-ui-system-folders-drives/1024/Developer-icon.png"
		item.LocationCity = []string{"上海", "北京", "成都"}
		item.ReviewCounts = rand.Int63n(100)
		item.Education = "本科"
		item.BusinessField = []string{"运维", "产品", "设计"}
		item.Days = rand.Intn(7)
		item.Months = rand.Intn(12)
		item.PayDay = rand.Intn(400)
		item.CanTransfer = i%2 == 0
		err := testDB.FirstOrCreate(&item).Error
		if err != nil {
			return
		}
	}
}

func createJobs() {

	for i := 0; i < 5; i++ {
		item := dbModel.CompuseJobs{}
		item.Id = "ffewfwef" + strconv.Itoa(i)
		item.Name = "实习大豆纤维多"
		item.CompanyID = "companyID"
		item.IconURL = "http://icons.iconarchive.com/icons/blackvariant/button-ui-system-folders-drives/1024/Developer-icon.png"
		item.LocationCity = []string{"city1", "city2", "city3"}
		item.ReviewCounts = 100
		item.Education = "大专"
		item.Type = "intern"

		err := testDB.FirstOrCreate(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func createOnlineApplys() {
	for i := 0; i < 3; i++ {
		item := dbModel.UserDeliveryStatusHistory{}

		t := time.Now().Add(time.Duration(rand.Int()) * time.Second)
		//t := time.Now().Add(time.Hour * 5)
		item.Status = i + 1
		item.Time = &t
		item.UserId = "1c0874f5-74c0-11e9-914b-a0999b089907"
		item.Describe = ""
		item.Type = "onlineApply"
		item.JobId = "1"
		err := testDB.Create(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func creatCarrerTalk() {
	for i := 0; i < 12; i++ {
		item := dbModel.CareerTalk{}
		item.Id = "宣讲会当前为多" + strconv.Itoa(i)
		item.Name = "宣讲会当前为多"
		item.CompanyID = "companyID"
		item.IconURL = "https://cdn1.iconfinder.com/data/icons/education-icons-3/155/vector_313_21-512.png"
		item.College = "我的大学"
		item.SimplifyAddress = "某个楼"
		t := time.Now()
		item.StartTime = &t
		item.ContentType = dbModel.Text

		err := testDB.Create(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createNews() {
	for i := 0; i < 6; i++ {
		item := dbModel.NewsModel{}
		var t = time.Now()
		item.Time = &t
		item.Icon = "http://pic34.photophoto.cn/20150128/0007020160374237_b.jpg"
		item.Link = "http://gorm.book.jasperxu.com/database.html"
		item.Title = "gorm 中文文档"
		item.Author = "github"

		testDB.Create(&item)
	}
}

func createRecruiter() {
	var before = time.Now().Add(-time.Hour)
	var r = dbModel.Recruiter{
		Uuid:      "12345",
		Name:      "我是小王",
		CompanyId: "公司1",
		Title:     "总监",
		Online:    true,
		UserIcon:  "http://www.sucaijishi.com/uploadfile/2018/0508/20180508023725429.png",
		LastLogin: &before,
	}
	err := testDB.Create(&r).Error
	if err != nil {
		print(err)
	}
}

func relatedQuery() {

	// 1
	var recruiters []dbModel.Recruiter
	var c dbModel.Company
	err := testDB.Model(&dbModel.Company{}).
		Where("id = ?", "companyID").Find(&c).Error

	err = testDB.Model(&c).Select("name").Association("Recruiters").Find(&recruiters).Error
	if err != nil {
		print(err.Error())
		return
	}
	fmt.Println(c)
	fmt.Println(recruiters)

	// 2

	var c1 dbModel.Company
	var fc = recruiters[0]
	// 默认使用主键id
	err = testDB.Model(fc).Related(&c1, "companyID").Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c1)

}

func testRecruitRelated() {

	//  3 r-c
	var rec1 dbModel.Recruiter
	err := testDB.Model(rec1).Where("uuid = ?", "12345").Find(&rec1).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rec1)

	err = testDB.Model(rec1).Association("CompusJobs").Find(&rec1.CompusJobs).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	err = testDB.Model(rec1).Association("InternJobs").Find(&rec1.InternJobs).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(rec1.CompusJobs), len(rec1.InternJobs))
}

func testVisitor() {
	testDB.Create(&dbModel.RecruiterVisitorUser{
		RecruiterId: "67f0c4ad-4888-11e9-b162-a0999b089907",
		UserId:      "24b12069-4853-11e9-a446-a0999b089907",
	})
}

func testSystemMessage() {
	var t = time.Now()
	//err := testDB.up(&dbModel.ForumReplyMyTime{
	//	LatestReplyTime: &t,
	//	UserId:          "24b12069-4853-11e9-a446-a0999b089907",
	//}).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	var mode dbModel.ForumThumbUpTime
	testDB.Model(&dbModel.ForumThumbUpTime{}).Where("user_id = ?", "24b12069-4853-11e9-a446-a0999b089907").First(&mode)
	testDB.Model(&mode).Update("latest_thumb_time", &t)
	fmt.Println(mode.CheckTime == nil, mode.LatestThumbTime)

}

func testForum() {

	// 创建帖子
	err := testDB.Create(&dbModel.ForumHotestArticle{
		Uuid: "2df1de31-6f16-11e9-8fa2-a0999b089907",
	}).Error
	if err != nil {
		fmt.Println(err)
		return
	}

}

func addUserDefaultTalk()  {
	var mes  = []string{"语句定位", "当前为多", "dqdwqfggg", "达瓦大"}

	err := testDB.Create(&dbModel.DefaultFirstMessage{
		Messages: mes,
		DefaultNum: 0,
		Open: true,
		UserId: "1c0874f5-74c0-11e9-914b-a0999b089907",
	}).Error
	if err != nil{
		fmt.Println(err)
	}

}

func main() {

	loadDB()
	//createJobs()
	//creatCarrerTalk()
	//createNews()
	//createOnlineApplys()
	//createInternJobs()
	//relatedQuery()
	//testRecruitRelated()
	//testVisitor()
	//testLeanCloud()
	//testSystemMessage()
	//testForum()
	addUserDefaultTalk()

}
