package sqlInfo

type RoomConfig struct {
	//ID       uint64 `db:"id" gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Svrid           uint32  `db:"Svrid" gorm:"column:Svrid;type:int(11);not null;default:0;" json:"Svrid,omitempty"`                                     //服务器id
	GameId          uint32  `db:"GameId" gorm:"column:GameId;type:int(11);not null;default:0;" json:"GameId,omitempty"`                                  //游戏ID
	RoomId          uint32  `db:"RoomId" gorm:"column:RoomId;type:int(11);not null;default:0;" json:"RoomId,omitempty"`                                  //房间id(svrid*1000*1000+gameid*1000+index)
	RoomIndex       uint32  `db:"RoomIndex" gorm:"column:RoomIndex;type:int(11);not null;default:0;" json:"RoomIndex,omitempty"`                         //房间index
	Base            uint32  `db:"Base" gorm:"column:Base;type:int(11);not null;default:0;" json:"Base,omitempty"`                                        //底注
	MinEntry        int     `db:"MinEntry" gorm:"column:MinEntry;type:int(11);not null;default:0;" json:"MinEntry"`                                      //进入限制(下) 0代表无限制
	MaxEntry        int     `db:"MaxEntry" gorm:"column:MaxEntry;type:int(11);not null;default:0;" json:"MaxEntry"`                                      //进入限制(上) 0代表无限制
	RoomName        string  `db:"RoomName" gorm:"column:RoomName;type:varchar(100);not null;default:'';" json:"RoomName,omitempty"`                      //房间名称
	RoomType        uint8   `db:"RoomType" gorm:"column:RoomType;type:int(11);not null;default:0;" json:"RoomType,omitempty"`                            //类型 1体验大厅 2正常大厅
	RoomSwitch      int     `db:"RoomSwitch" gorm:"column:RoomSwitch;type:tinyint(4);not null;default:0;" json:"RoomSwitch,omitempty"`                   //房间开关
	RoomWelfare     int     `db:"RoomWelfare" gorm:"column:RoomWelfare;type:tinyint(4);not null;default:0;" json:"RoomWelfare,omitempty"`                //房间赠送
	RoomWelfareDesc *string `db:"RoomWelfareDesc" gorm:"column:RoomWelfareDesc;type:varchar(255);not null;default:'';" json:"RoomWelfareDesc,omitempty"` //房间赠送描述
	Desc            *string `db:"Desc" gorm:"column:Desc;type:varchar(255);not null;default:'';" json:"Desc,omitempty"`                                  //房间描述
	Tax             int     `db:"Tax" gorm:"column:Tax;type:int(11);not null;default:0;" json:"Tax,omitempty"`                                           //税千分比
	BonusDiscount   int     `db:"BonusDiscount" gorm:"column:BonusDiscount;type:int(11);not null;default:0;" json:"BonusDiscount,omitempty"`             //比例千分比
	AiSwitch        int     `db:"AiSwitch" gorm:"column:AiSwitch;type:tinyint(4);not null;default:0;" json:"AiSwitch,omitempty"`                         //ai开关
	AiLimit         int     `db:"AiLimit" gorm:"column:AiLimit;type:int(11);not null;default:0;" json:"AiLimit,omitempty"`                               //ai人数限制
	ExtData         string  `db:"ExtData" gorm:"column:ExtData;type:varchar(500);not null;default:'';" json:"ExtData,omitempty"`                         //特殊配置
	PoolID          int     `db:"PoolID" gorm:"column:PoolID;type:int(11);not null;default:0;" json:"PoolID,omitempty"`
	RechargeLimit   int     `db:"RechargeLimit" gorm:"column:RechargeLimit;type:int(11);not null;default:0;" json:"RechargeLimit"`
	PExtData        *string `db:"pExtData" gorm:"column:pExtData;type:varchar(500);not null;default:'';" json:"pExtData,omitempty"` //奖池配置
	Cnt             any     `gorm:"column:cnt;type:varchar(500);not null;default:'';" json:"cnt,omitempty"`
	RechargeCfg     any     `json:"RechargeCfg,omitempty"`
	RechargeCount   int     `db:"RechargeCount" json:"RechargeCount,omitempty"`
	CurPoolValue    int64   `db:"CurPoolValue" json:"CurPoolValue,omitempty"`
	OnlineMin       int32   `db:"OnlineMin" json:"-"` //最低随机在线人数
	OnlineMax       int32   `db:"OnlineMax" json:"-"` //最高随机在线人数
}

func (table *RoomConfig) TableName() (tableName string) {
	tableName = "roomview" //视图
	return
}

// 充值礼包配置
type RechargePackConfig struct {
	ID           uint64 `db:"id"  json:"id,omitempty"`
	GameId       int    `db:"game_id" json:"game_id,omitempty"`              //游戏id
	RoomId       int    `db:"room_id"  json:"room_id,omitempty"`             //房间id
	BasicRewards int    `db:"basic_rewards"  json:"basic_rewards,omitempty"` //基础奖励
	BasicType    uint8  `db:"basic_type"  json:"basic_type,omitempty"`       //金额类别(1=可提现可下注,2=不可提现可下注)
	GiftRewards  int    `db:"gift_rewards"  json:"gift_rewards,omitempty"`   // 赠送额度
	GiftType     uint8  `db:"gift_type"  json:"gift_type,omitempty"`         //金额类别(1=可提现可下注,2=不可提现可下注)
	Bonus        int    `db:"bonus"  json:"bonus,omitempty"`                 //bonus
	BonusType    uint8  `db:"bonus_type"  json:"bonus_type,omitempty"`       //金额类别(1=可提现可下注,2=不可提现可下注)
	Total        int    `db:"total"  json:"total,omitempty"`                 //总额度
	Price        int    `db:"price" json:"price,omitempty"`                  //价格
	Times        int    `db:"times"  json:"times,omitempty"`                 //出现次数
	Interval     int    `db:"interval"  json:"interval,omitempty"`           //间隔时长(s)
	Ratio        int    `db:"ratio"  json:"ratio,omitempty"`                 //折扣比例
}

func (table *RechargePackConfig) TableName() (tableName string) {
	return "recharge_pack_config"
}
