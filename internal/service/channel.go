// internal/service/channel.go
package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/example/epay-go/internal/model"
	"github.com/example/epay-go/internal/repository"
	"github.com/shopspring/decimal"
)

type ChannelService struct {
	repo *repository.ChannelRepository
}

func NewChannelService() *ChannelService {
	return &ChannelService{
		repo: repository.NewChannelRepository(),
	}
}

// CreateChannelRequest 创建通道请求
type CreateChannelRequest struct {
	Name       string                 `json:"name" binding:"required"`
	Key        string                 `json:"key"`
	Plugin     string                 `json:"plugin" binding:"required"`
	PayTypes   string                 `json:"pay_types"`
	AppType    string                 `json:"app_type"` // 支持的接口类型
	Config     map[string]interface{} `json:"config"`
	Rate       float64                `json:"rate"`
	DailyLimit float64                `json:"daily_limit"`
	Status     int8                   `json:"status"`
	Sort       int                    `json:"sort"`
}

// Create 创建通道
func (s *ChannelService) Create(req *CreateChannelRequest) (*model.Channel, error) {
	// 验证 AppType 不能为空
	if req.AppType == "" {
		return nil, fmt.Errorf("请至少选择一个支付接口")
	}

	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}

	channel := &model.Channel{
		Name:       req.Name,
		Key:        resolveChannelKey(req),
		Plugin:     req.Plugin,
		PayTypes:   req.PayTypes,
		AppType:    req.AppType,
		Config:     configJSON,
		Rate:       decimal.NewFromFloat(req.Rate),
		DailyLimit: decimal.NewFromFloat(req.DailyLimit),
		Status:     req.Status,
		Sort:       req.Sort,
	}

	if err := s.repo.Create(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func resolveChannelKey(req *CreateChannelRequest) string {
	if req == nil {
		return "channel_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	trimmed := strings.TrimSpace(req.Key)
	if trimmed != "" {
		return trimmed
	}
	plugin := strings.TrimSpace(req.Plugin)
	if plugin == "" {
		plugin = "channel"
	}
	return plugin + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// GetByID 根据ID获取通道
func (s *ChannelService) GetByID(id int64) (*model.Channel, error) {
	return s.repo.GetByID(id)
}

// UpdateChannelRequest 更新通道请求
type UpdateChannelRequest struct {
	Name       string                 `json:"name"`
	PayTypes   string                 `json:"pay_types"`
	AppType    string                 `json:"app_type"`
	Config     map[string]interface{} `json:"config"`
	Rate       float64                `json:"rate"`
	DailyLimit float64                `json:"daily_limit"`
	Status     int8                   `json:"status"`
	Sort       int                    `json:"sort"`
}

// Update 更新通道
func (s *ChannelService) Update(id int64, req *UpdateChannelRequest) error {
	channel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		channel.Name = req.Name
	}
	if req.PayTypes != "" {
		channel.PayTypes = req.PayTypes
	}
	if req.AppType != "" {
		channel.AppType = req.AppType
	}
	if req.Config != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			return err
		}
		channel.Config = configJSON
	}
	channel.Rate = decimal.NewFromFloat(req.Rate)
	channel.DailyLimit = decimal.NewFromFloat(req.DailyLimit)
	channel.Status = req.Status
	channel.Sort = req.Sort

	return s.repo.Update(channel)
}

// Delete 删除通道
func (s *ChannelService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// List 分页查询通道列表
func (s *ChannelService) List(page, pageSize int) ([]model.Channel, int64, error) {
	return s.repo.List(page, pageSize)
}

// ListEnabled 获取所有启用的通道
func (s *ChannelService) ListEnabled() ([]model.Channel, error) {
	return s.repo.ListEnabled()
}

// GetAvailableChannel 根据支付类型获取可用通道
func (s *ChannelService) GetAvailableChannel(payType string) (*model.Channel, error) {
	return s.repo.GetAvailableByPayType(payType)
}
