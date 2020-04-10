package Steam

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CommunityVisibilityState int

const (
	// Private community visibility state
	Private CommunityVisibilityState = 1
	// FriendsOnly community visibility state
	FriendsOnly CommunityVisibilityState = 2
	// Public community visibility state
	Public CommunityVisibilityState = 3
)

var VisibilityState = []string{"", "隐私", "好友可见", "公开"}

type Profile struct {
	OnlineState     string `xml:"onlineState"`
	AvatarIcon      string `xml:"avatarIcon"`
	VisibilityState int    `xml:"visibilityState"`
}

func GetProfile(steamID64 string) (*Profile, error) {
	url := fmt.Sprintf("http://steamcommunity.com/profiles/%s/?xml=1", steamID64)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	profile := Profile{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//a := "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><profile><steamID64>76561198027072210</steamID64><steamID><![CDATA[VAN.CY]]></steamID><onlineState>online</onlineState><stateMessage><![CDATA[Online]]></stateMessage><privacyState>public</privacyState><visibilityState>3</visibilityState><avatarIcon><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/96/965b48a8673e0e684c1790380dd1b842430343c5.jpg]]></avatarIcon><avatarMedium><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/96/965b48a8673e0e684c1790380dd1b842430343c5_medium.jpg]]></avatarMedium><avatarFull><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/96/965b48a8673e0e684c1790380dd1b842430343c5_full.jpg]]></avatarFull><vacBanned>0</vacBanned><tradeBanState>None</tradeBanState><isLimitedAccount>0</isLimitedAccount><customURL><![CDATA[VanSi]]></customURL><memberSince>June 27, 2010</memberSince><steamRating></steamRating><hoursPlayed2Wk>0.0</hoursPlayed2Wk><headline><![CDATA[]]></headline><location><![CDATA[]]></location><realname><![CDATA[]]></realname><summary><![CDATA[<br><br><br><br><br>It takes only a minute to get a crush on someone, an hour to like someone, and a day to love someone, but it takes a lifetime to forget someone. <br>-<br><br><br>你做了选择，对的错的，我只能承认，心是痛的。怀疑你说的，我被伤的那麽深，就放声哭了，何必再强忍。<br>我没有选择，我不再完整，你只能默认，我要被割舍。<br>如果这不是结局，如果我还爱你。如果我愿相信，你就是唯一。如果你听到这里，如果你依然放弃，那这就是爱情，我难以抗拒。<br>如果这就是爱情，本来就不公平。你不需要讲理，我可以离去。如果我成全了你，如果我能祝福你，那不是我看清，是我证明，我爱你。<br><br>灰色的天空，无法猜透，多余的眼泪，无法挽留。<br>什么都牵动，我感觉真的好脆弱。<br>我不要你走，我不想放手，却又不能够奢求，同情的温柔。<br>你可以自由，我愿意承受，把昨天留给我。]]></summary><groups><group isPrimary=\"1\"><groupID64>103582791454977334</groupID64><groupName><![CDATA[Coweye]]></groupName><groupURL><![CDATA[coweye]]></groupURL><headline><![CDATA[Are you a kawayeh ugu qt 7/10 azn girl for daddy? Yes you are gorl!]]></headline><summary><![CDATA[This is a group solely for azn anime gorls with a burning passion for daddies’ cummies and anal rape with Gucci chokers. If you’re man enough to become a sissy cuck crossdresser like me, please join this group and invite your friends, because prostate autism is the only cure against cancer and dysfunctional testies.<br><br>YEAH!<br>Let's do it altogether now!<br>YEAH!<br>*bawl of a thousand masochistic manchildren*<br>NIKO NIKO NIIIIIIIIIII~]]></summary><avatarIcon><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fc/fcf892b9bc8cd87c85249388176ca0f9dbc17918.jpg]]></avatarIcon><avatarMedium><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fc/fcf892b9bc8cd87c85249388176ca0f9dbc17918_medium.jpg]]></avatarMedium><avatarFull><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fc/fcf892b9bc8cd87c85249388176ca0f9dbc17918_full.jpg]]></avatarFull><memberCount>255</memberCount><membersInChat>2</membersInChat><membersInGame>10</membersInGame><membersOnline>57</membersOnline></group><group isPrimary=\"0\"><groupID64>103582791429538047</groupID64><groupName><![CDATA[Xtreme-Jumps]]></groupName><groupURL><![CDATA[Xtreme-Jumps]]></groupURL><headline><![CDATA[]]></headline><summary><![CDATA[No information given.]]></summary><avatarIcon><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/a7/a738f173d6f774b3209355e6951ddd2a6777d62c.jpg]]></avatarIcon><avatarMedium><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/a7/a738f173d6f774b3209355e6951ddd2a6777d62c_medium.jpg]]></avatarMedium><avatarFull><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/a7/a738f173d6f774b3209355e6951ddd2a6777d62c_full.jpg]]></avatarFull><memberCount>9738</memberCount><membersInChat>118</membersInChat><membersInGame>262</membersInGame><membersOnline>2038</membersOnline></group><group isPrimary=\"0\"><groupID64>103582791429697786</groupID64><groupName><![CDATA[97club]]></groupName><groupURL><![CDATA[97Club]]></groupURL><headline><![CDATA[]]></headline><summary><![CDATA[No information given.]]></summary><avatarIcon><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/58/58279c06f0144f61550d00be2b173053908d6241.jpg]]></avatarIcon><avatarMedium><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/58/58279c06f0144f61550d00be2b173053908d6241_medium.jpg]]></avatarMedium><avatarFull><![CDATA[https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/58/58279c06f0144f61550d00be2b173053908d6241_full.jpg]]></avatarFull><memberCount>1438</memberCount><membersInChat>20</membersInChat><membersInGame>25</membersInGame><membersOnline>176</membersOnline></group><group isPrimary=\"0\"><groupID64>103582791429732618</groupID64></group><group isPrimary=\"0\"><groupID64>103582791430690152</groupID64></group><group isPrimary=\"0\"><groupID64>103582791430760692</groupID64></group><group isPrimary=\"0\"><groupID64>103582791432145822</groupID64></group><group isPrimary=\"0\"><groupID64>103582791432229965</groupID64></group><group isPrimary=\"0\"><groupID64>103582791432857211</groupID64></group><group isPrimary=\"0\"><groupID64>103582791433376641</groupID64></group><group isPrimary=\"0\"><groupID64>103582791433837621</groupID64></group><group isPrimary=\"0\"><groupID64>103582791434208261</groupID64></group><group isPrimary=\"0\"><groupID64>103582791436563633</groupID64></group><group isPrimary=\"0\"><groupID64>103582791437140637</groupID64></group></groups></profile>"
	err = xml.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
