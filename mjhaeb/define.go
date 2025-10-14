package mjhaeb

const (
	RuleDiscardTime    = iota //出牌时间 0
	RuleWaitTime              //等待时间 1
	RuleScoreType             //算分方式 2
	RuleHuType                //胡牌类型 3 (0-任意胡，1-死卡)
	RuleLouBao                //搂宝 4 （1-打开，0-关闭）
	RuleBaoZhongBao           //宝中宝 5 (1-打开，0-关闭)
	RuleHZMTF                 //红中满天飞 6 (1-打开，0-关闭)
	RuleGuaDaFeng             //刮大风 7 (1-打开，0-关闭)
	RuleTingChuShouPao        //报听出手炮 8 （1-打开，0-关闭）
	Rule37Jia                 //3、7夹 9 (1-打开，0-关闭)
	RuleDanDiaoJia            //单吊夹 10 (1-打开，0-关闭)
	RuleDuiPengJia            //对碰夹 11 (1-打开，0-关闭)
	RuleDuoMianTingJia        //多面听夹 12 (1-打开，0-关闭)
	Rule7Dui                  //七对 13 (1-打开，0-关闭)
	RuleKonSouFen             //杠收分 14 (1-打开，0-关闭)
	RuleQinYiSe               //清一色 15 (1-打开，0-关闭)
	RuleMengQin3Bei           //门清*3 16 (1-打开，0-关闭)
	RuleXianDaHouMo           //先打后摸 17 (1-打开，0-关闭)
	RuleEnd
)
