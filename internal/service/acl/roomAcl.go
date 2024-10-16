package acl

import (
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef/enum"

	"github.com/mohae/deepcopy"
)

type RoomAcl struct {
	UserId       string
	TargetUserId string
	RoomId       string
	Seat         int
	Scene        string
	roomInfo     model.Room
}

func (r *RoomAcl) GetAcl() any {
	roomInfo, _ := new(dao.RoomDao).GetRoomById(r.RoomId)
	if len(roomInfo.Id) > 0 {
		r.roomInfo = roomInfo
	}
	var res any
	switch r.Scene {
	case "mic": //触发场景麦位
		res = r.getMicAcl()
	case "card": //资料卡
		micInfo := r.GetUserMicInfo(r.RoomId, r.TargetUserId)
		if micInfo != nil {
			r.Seat = micInfo.Id
			return r.getMicAcl()
		}
		res = r.getUserCardAcl()
	case "more": //更多选项
		res = r.getRoomMoreAcl()
	case "hidden_mic":
		res = r.getHiddenMicAcl()
	}
	return res
}

func (r *RoomAcl) getHiddenMicAcl() []Menu {
	rules := []string{}
	onMic := r.IsOnHiddenMic(r.RoomId, r.UserId)
	if onMic {
		rules = append(rules, DownHiddenMic)
	} else {
		rules = append(rules, UpHiddenMic)
	}
	authMap, err := r.getVerifyAuth(r.UserId, r.RoomId, rules...)
	if err != nil {
		panic(err)
	}
	res := []Menu{}
	for _, v := range HiddenMic {
		if has, ok := authMap[v.Name]; ok && has {
			clearMenu := v
			clearMenu.Show = true
			clearMenu = r.clearRules(clearMenu)
			clearMenu = r.clearAnchorRules(clearMenu)
			if clearMenu.Show {
				res = append(res, clearMenu)
			}
		}
	}
	return res
}

func (r *RoomAcl) getMicAcl() any {
	info := r.MicUserInfo(r.RoomId, r.Seat)
	//判断用户是否拥有关闭麦位权限
	if info.Status == enum.MicStatusClose {
		switch info.Identity {
		case enum.CompereMicSeat:
			ok, _ := r.CheckUserRule(r.UserId, r.RoomId, SwitchCompereMic)
			if !ok {
				panic(error2.I18nError{
					Code: error2.ErrCodeRoomSeatClosed,
				})
			}
		case enum.GuestMicSeat:
			ok, _ := r.CheckUserRule(r.UserId, r.RoomId, SwitchGuestMic)
			if !ok {
				panic(error2.I18nError{
					Code: error2.ErrCodeRoomSeatClosed,
				})
			}
		case enum.CounselorMicSeat:
			ok, _ := r.CheckUserRule(r.UserId, r.RoomId, SwitchCounselorMic)
			if !ok {
				panic(error2.I18nError{
					Code: error2.ErrCodeRoomSeatClosed,
				})
			}
		case enum.MusicianMicSeat:
			ok, _ := r.CheckUserRule(r.UserId, r.RoomId, SwitchMusicianMic)
			if !ok {
				panic(error2.I18nError{
					Code: error2.ErrCodeRoomSeatClosed,
				})
			}
		case enum.NormalMicSeat:
			ok, _ := r.CheckUserRule(r.UserId, r.RoomId, SwitchNormalMic)
			if !ok {
				panic(error2.I18nError{
					Code: error2.ErrCodeRoomSeatClosed,
				})
			}
		}
	}
	if len(info.UserInfo.UserId) == 0 { //麦位没有人
		switch info.Identity {
		case enum.CompereMicSeat: //主持麦
			return r.getEmptyMicAcl(ZhuChiEmptyMic)
		case enum.GuestMicSeat: //嘉宾麦
			return r.getEmptyMicAcl(JiaBinEmptyMic)
		case enum.MusicianMicSeat: //音乐人麦
			return r.getEmptyMicAcl(YinYueRenEmptyMic)
		case enum.CounselorMicSeat: //咨询师麦
			return r.getEmptyMicAcl(ZiXunShiEmptyMic)
		case enum.NormalMicSeat: //普通麦
			return r.getEmptyMicAcl(NormalEmptyMic)
		}
	} else { //麦位有人
		switch info.Identity {
		case enum.CompereMicSeat: //主持麦
			return r.getUserMicAcl(ZhuChiUserMic)
		case enum.GuestMicSeat: //嘉宾麦
			return r.getUserMicAcl(JiaBinUserMic)
		case enum.MusicianMicSeat: //音乐人麦
			return r.getUserMicAcl(YinYueRenUserMic)
		case enum.CounselorMicSeat: //咨询师麦
			return r.getUserMicAcl(ZiXunShiUserMic)
		case enum.NormalMicSeat: //普通麦
			return r.getUserMicAcl(NormalUserMic)
		}
	}

	return nil
}

// 麦位有用户权限表
func (r *RoomAcl) getUserMicAcl(acl []Menu) any {
	rules := []string{}
	for _, v := range acl {
		rules = append(rules, v.Name)
	}
	authMap, err := r.getVerifyAuth(r.UserId, r.RoomId, rules...)
	if err != nil {
		panic(err)
	}
	res := []Menu{}
	for _, v := range acl {
		if has, ok := authMap[v.Name]; ok && has {
			clearMenu := v
			clearMenu.Show = true
			clearMenu = r.clearRules(clearMenu)
			clearMenu = r.clearAnchorRules(clearMenu)
			if clearMenu.Show {
				clearMenu = r.getSwitch(clearMenu)
				res = append(res, clearMenu)
			}
		}
	}
	return res
}

// 空麦位权限表
func (r *RoomAcl) getEmptyMicAcl(acl []Menu) []Menu {
	rules := []string{}
	for _, v := range acl {
		rules = append(rules, v.Name)
	}
	authMap, err := r.getVerifyAuth(r.UserId, r.RoomId, rules...)
	if err != nil {
		panic(err)
	}
	//r.excludeCompereRule(r.UserId, r.RoomId, authMap)
	res := []Menu{}
	for _, v := range acl {
		if has, ok := authMap[v.Name]; ok && has {
			clearMenu := v
			clearMenu.Show = true
			clearMenu = r.clearRules(clearMenu)
			clearMenu = r.clearAnchorRules(clearMenu)
			if clearMenu.Show {
				clearMenu = r.getSwitch(clearMenu)
				res = append(res, clearMenu)
			}
		}
	}
	return res
}

// 用户资料卡权限表
func (r *RoomAcl) getUserCardAcl() []Menu {
	rules := []string{}
	for _, v := range UserCard {
		rules = append(rules, v.Name)
	}
	authMap, err := r.getVerifyAuth(r.UserId, r.RoomId, rules...)
	if err != nil {
		panic(err)
	}
	res := []Menu{}
	for _, v := range UserCard {
		if has, ok := authMap[v.Name]; ok && has {
			clearMenu := v
			clearMenu.Show = true
			clearMenu = r.clearRules(v, false)
			clearMenu = r.clearAnchorRules(clearMenu)
			if clearMenu.Show {
				clearMenu = r.getSwitch(clearMenu)
				res = append(res, clearMenu)
			}
		}
	}
	return res
}

// 获取更多权限列表
func (r *RoomAcl) getRoomMoreAcl() map[string][]Menu {
	rules := []string{}
	for _, v := range HuDong {
		rules = append(rules, v.Name)
	}
	for _, v := range Manage {
		rules = append(rules, v.Name)
	}
	for _, v := range Other {
		rules = append(rules, v.Name)
	}
	authMap, err := r.getVerifyAuth(r.UserId, r.RoomId, rules...)
	if err != nil {
		panic(err)
	}
	data := deepcopy.Copy(MoreMenu).(map[string][]Menu)
	for key, value := range data {
		for k, menu := range value {
			if has, ok := authMap[menu.Name]; ok && has {
				clearMenu := menu
				clearMenu.Show = true
				clearMenu = r.clearRules(clearMenu)
				clearMenu = r.clearAnchorRules(clearMenu)
				if clearMenu.Show {
					clearMenu = r.getSwitch(clearMenu)
					data[key][k] = clearMenu
				}
			}
		}
	}
	res := map[string][]Menu{}
	for key, value := range data {
		menus := []Menu{}
		for _, menu := range value {
			if menu.Show {
				menus = append(menus, menu)
			}
		}
		res[key] = menus
	}
	return res
}

// 处理状态类开关类权限
func (r *RoomAcl) getSwitch(menu Menu) Menu {
	switch menu.Name {
	case HiddenRoom: //隐藏房间开启关闭状态
		//判断用户是否开启隐藏房间
		status := r.RoomHiddenStatus(r.RoomId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启房间"
			return menu
		}
	case LockRoom: //锁定房间开启关闭状态
		//判断是否锁定房间
		status := r.RoomLockStatus(r.RoomId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "解锁房间"
			return menu
		}
	case FreedMic: //自由上下麦开启关闭状态
		//判断房间是否开启自由上下麦
		status := r.RoomFreedMicStatus(r.RoomId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "关闭自由上下麦"
			return menu
		}
	case FreedSpeak: //自由发言开启关闭状态
		status := r.RoomFreedSpeakStatus(r.RoomId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "关闭自由发言"
			return menu
		}
	case RoomMute: //房间静音
		status := r.RoomMuteStatus(r.RoomId, r.UserId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "关闭静音"
			return menu
		}
	case RoomCloseSpecialEffects: //关闭动效开启关闭状态
		status := r.RoomSpecialEffectsStatus(r.RoomId, r.UserId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启动效"
			return menu
		}
	case RoomClosePublicChat: //关闭公屏开启关闭状态
		status := r.RoomPublicChatStatus(r.RoomId)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "关闭公屏"
			return menu
		}
	case SwitchCompereMic: //关闭主持麦位开启关闭状态
		status := r.MicSwitchStatus(r.RoomId, r.Seat, enum.CompereMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启麦位"
			return menu
		}
	case SwitchGuestMic: //关闭嘉宾麦位开启关闭状态
		status := r.MicSwitchStatus(r.RoomId, r.Seat, enum.GuestMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启麦位"
			return menu
		}
	case SwitchMusicianMic: //关闭音乐人麦位开启关闭状态
		status := r.MicSwitchStatus(r.RoomId, r.Seat, enum.MusicianMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启麦位"
			return menu
		}
	case SwitchCounselorMic: //关闭咨询师麦位开启关闭状态
		status := r.MicSwitchStatus(r.RoomId, r.Seat, enum.CounselorMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启麦位"
			return menu
		}
	case SwitchNormalMic: //关闭普通麦位开启关闭状态
		status := r.MicSwitchStatus(r.RoomId, r.Seat, enum.NormalMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "开启麦位"
			return menu
		}
	case MuteCompereMic: //主持麦静音
		status := r.MicMuteStatus(r.RoomId, r.Seat, enum.CompereMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "取消静音"
			return menu
		}
	case MuteGuestMic: //嘉宾麦静音
		status := r.MicMuteStatus(r.RoomId, r.Seat, enum.GuestMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "取消静音"
			return menu
		}
	case MuteNormalMic: //普通麦静音
		status := r.MicMuteStatus(r.RoomId, r.Seat, enum.NormalMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "取消静音"
			return menu
		}
	case MuteMusicianMic: //音乐人麦静音
		status := r.MicMuteStatus(r.RoomId, r.Seat, enum.MusicianMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "取消静音"
			return menu
		}
	case MuteCounselorMic: //咨询师静音
		status := r.MicMuteStatus(r.RoomId, r.Seat, enum.CounselorMicSeat)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "取消静音"
			return menu
		}
	case RoomRelateWheat: //连麦操作
		status := r.RoomRelateWheat(r.roomInfo)
		if status == enum.SwitchOpen {
			menu.Switch = status
			menu.Title = "关闭连麦"
		} else {
			menu.Switch = enum.SwitchOff
			menu.Title = "开启连麦"
		}
		return menu
	}
	return menu
}

func (r *RoomAcl) RoomRelateWheat(roomInfo model.Room) int {
	//判断个播房间是否在连麦
	if roomInfo.LiveType == enum.LiveTypeAnchor {
		if _, ok := enum.RoomTemplates[roomInfo.TemplateId]; ok {
			return enum.SwitchOpen
		}
	}
	return enum.SwitchOff
}

// 对规则进行清洗处理
func (r *RoomAcl) clearRules(menu Menu, onMic ...bool) Menu {
	onSeat := true
	if len(onMic) > 0 {
		onSeat = onMic[0]
	}
	switch menu.Name {
	case DownCompereMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = true
		} else {
			menu.Show = false
		}
	case DownGuestMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = true
		} else {
			menu.Show = false
		}
	case DownMusicianMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = true
		} else {
			menu.Show = false
		}
	case DownCounselorMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = true
		} else {
			menu.Show = false
		}
	case DownNormalMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = true
		} else {
			menu.Show = false
		}
	case HoldUserDownNormalMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = false
		} else {
			menu.Show = true
		}
	case UpNormalMic: //个播间普通麦位，房主点击普通麦位，需要屏蔽上麦权限
		if r.roomInfo.UserId == r.UserId && r.roomInfo.LiveType == enum.LiveTypeAnchor {
			menu.Show = false
		}
	case ReportUser:
		if onSeat {
			info := r.MicUserInfo(r.RoomId, r.Seat)
			if info.UserInfo.UserId == r.UserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		} else {
			if r.UserId == r.TargetUserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		}

	case HoldCompereDownCompereMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = false
		} else {
			menu.Show = true
		}
	case HoldUserDownGuestMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = false
		} else {
			menu.Show = true
		}
	case HoldUserDownMusicianMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = false
		} else {
			menu.Show = true
		}
	case HoldUserDownCounselorMic:
		info := r.MicUserInfo(r.RoomId, r.Seat)
		if info.UserInfo.UserId == r.UserId {
			menu.Show = false
		} else {
			menu.Show = true
		}
	case OutRoom:
		if onSeat {
			info := r.MicUserInfo(r.RoomId, r.Seat)
			if info.UserInfo.UserId == r.UserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		} else {
			if r.UserId == r.TargetUserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		}
	case AddRoomBlacklist:
		if onSeat {
			info := r.MicUserInfo(r.RoomId, r.Seat)
			if info.UserInfo.UserId == r.UserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		} else {
			if r.UserId == r.TargetUserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		}
	case ShutUp:
		if onSeat {
			info := r.MicUserInfo(r.RoomId, r.Seat)
			if info.UserInfo.UserId == r.UserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		} else {
			if r.UserId == r.TargetUserId {
				menu.Show = false
			} else {
				menu.Show = true
			}
		}
	}
	return menu
}

// 清理个播房规则
func (r *RoomAcl) clearAnchorRules(menu Menu) Menu {
	if r.roomInfo.RoomType != enum.RoomTypeAnchorVoice {
		return menu
	}
	switch menu.Name {
	case DownCompereMic:
		menu.Show = false
	case HoldCompereDownCompereMic:
		menu.Show = false
	case UpCompereMic, SwitchCompereMic, HoldCompereUpCompereMic:
		menu.Show = false
	case UpHiddenMic, DownHiddenMic:
		if r.roomInfo.UserId == r.UserId {
			menu.Show = false
		}
	case FreedMic, FreedSpeak, RoomClearMic, MuteCompereMic:
		menu.Show = false
	}
	return menu
}
