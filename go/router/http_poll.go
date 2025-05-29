package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// OnePoll 获取所有投票中的单个投票选项
type OnePoll struct {
	OptionId          int    `json:"option_id"`          // 投票选项1,2,3,4
	OptionDescription string `json:"option_description"` // 选项描述
	VoteCount         int    `json:"vote_count"`         // 选项票数
}

// AllPolls 所有投票结构体
type AllPolls struct {
	Polls []OnePoll `json:"polls"` // 所有投票选项
}

// VoteRequest 用户投票请求结构体
type VoteRequest struct {
	UserUUID   string `json:"user_uuid" binding:"required"`
	VoteOption int    `json:"vote_option" binding:"required,min=1,max=4"`
}

func InitPollHttpRouter(routerGroup *gin.RouterGroup, db *gorm.DB) {
	pollRouter := routerGroup.Group("/poll")
	{
		pollRouter.GET("/", func(c *gin.Context) {
			GetAllPolls(c, db)
		}) // 获取所有投票
		pollRouter.POST("/vote", func(c *gin.Context) {
			PostPoll(c, db)
		}) // 提交投票
	}
}

func GetAllPolls(c *gin.Context, db *gorm.DB) {
	rows, err := db.Raw("SELECT option_id, option_description, vote_count FROM vote_statistics ORDER BY option_id").Rows()
	if err != nil {
		log.Printf("查询投票数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票数据失败"})
		return
	}
	defer rows.Close()

	var polls []OnePoll
	for rows.Next() {
		var poll OnePoll
		if err := db.ScanRows(rows, &poll); err != nil {
			log.Printf("扫描行数据失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "处理投票数据失败"})
			return
		}
		polls = append(polls, poll)
	}

	c.JSON(http.StatusOK, AllPolls{Polls: polls})
}

func PostPoll(c *gin.Context, db *gorm.DB) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开始事务失败"})
		return
	}

	// 添加用户投票记录
	if err := tx.Exec("INSERT INTO user_votes (user_uuid, vote_option) VALUES (?, ?)",
		req.UserUUID, req.VoteOption).Error; err != nil {
		tx.Rollback()
		log.Printf("添加用户投票记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "投票失败"})
		return
	}

	// 更新投票统计
	if err := tx.Exec("UPDATE vote_statistics SET vote_count = vote_count + 1 WHERE option_id = ?",
		req.VoteOption).Error; err != nil {
		tx.Rollback()
		log.Printf("更新投票统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票统计失败"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "投票处理失败"})
		return
	}

	go broadcastPollUpdate(db)

	c.JSON(http.StatusOK, gin.H{"message": "投票成功"})
}

func broadcastPollUpdate(db *gorm.DB) {
	rows, err := db.Raw("SELECT option_id, option_description, vote_count FROM vote_statistics ORDER BY option_id").Rows()
	if err != nil {
		log.Printf("查询投票数据失败: %v", err)
		return
	}
	defer rows.Close()

	var polls []OnePoll
	for rows.Next() {
		var poll OnePoll
		if err := db.ScanRows(rows, &poll); err != nil {
			log.Printf("扫描行数据失败: %v", err)
			return
		}
		polls = append(polls, poll)
	}

	pollData, err := json.Marshal(AllPolls{Polls: polls})
	if err != nil {
		log.Printf("序列化投票数据失败: %v", err)
		return
	}

	wsManager.BroadcastPollUpdate(pollData)
}
