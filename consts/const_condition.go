package consts

/**该文件，存储所有条件id*/

/**弹窗触发条件（未达到条件则不弹窗）*/
const (
	ConditonNo                int = 0 //无条件
	ConditonCoin              int = 1 //玩家的金额低于?时
	ConditonRechargeSmall     int = 2 //玩家充值总额低于?时
	ConditonRechargeBigger    int = 3 //玩家充值总额高于？时
	ConditonSendRateSmall     int = 4 //玩家赠送比低于 ? 时
	ConditonSendRateBigger    int = 5 //玩家赠送比高于 ? 时
	ConditonCoinSmallRecharge int = 6 //玩家总金币小于总充值比例 ?
)
