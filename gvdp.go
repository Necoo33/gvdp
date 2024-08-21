package gvdp

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type VDPUser struct {
	Ip             string
	Count          uint16
	Banned         bool
	DateToBeBanned int64
	StartTime      int64
}

type VanillaDdosProtector struct {
	IndividualUsers []VDPUser
	AttackTimespan  uint16
	AttackCount     uint16
	BanTime         uint
	BannedUsers     []VDPUser
	ErrorCode       int
	Whitelist       []string
	BanOccured      bool
}

func Init(at uint16, ac uint16, bt uint, ec int, wl []string) VanillaDdosProtector {
	return VanillaDdosProtector{
		IndividualUsers: []VDPUser{
			{
				Ip:             "999.999.999.999",
				Count:          1,
				Banned:         false,
				DateToBeBanned: -1,
				StartTime:      -1,
			},
		},
		AttackTimespan: at,
		AttackCount:    ac,
		BanTime:        bt,
		BannedUsers: []VDPUser{
			{
				Ip:             "999.999.999",
				Count:          1,
				Banned:         true,
				DateToBeBanned: 10000000000000,
				StartTime:      -1,
			},
		},
		ErrorCode: ec,
		Whitelist: wl,
	}
}

func (vdp VanillaDdosProtector) HandleBanningAndAllowing(req *http.Request) VanillaDdosProtector {
	repeatedAttack := false
	banEnded := false
	userWhichReallowed := VDPUser{
		Ip:             "234234",
		Count:          2342,
		Banned:         false,
		DateToBeBanned: -1,
		StartTime:      -1,
	}

	for index, iu := range vdp.IndividualUsers {
		if iu.Banned {
			now := time.Now().Unix()

			if iu.DateToBeBanned != -1 && now > iu.DateToBeBanned {
				vdp.IndividualUsers[index].Banned = false

				banEnded = true

				vdp.IndividualUsers[index].DateToBeBanned = -1

				vdp.IndividualUsers[index].StartTime = -1

				vdp.IndividualUsers[index].Count = 0

				userWhichReallowed.Ip = vdp.IndividualUsers[index].Ip

				userWhichReallowed.Count = 0
			}
		}
	}

	if userWhichReallowed.Count != 2342 {
		for index, iu := range vdp.IndividualUsers {
			if iu.Ip == userWhichReallowed.Ip {
				vdp.IndividualUsers = append(vdp.IndividualUsers[:index], vdp.IndividualUsers[index+1:]...)

				vdp.BanOccured = false
			}
		}
	}

	clientIp, _, _ := net.SplitHostPort(req.RemoteAddr)

	if !banEnded {
		for index, iu := range vdp.IndividualUsers {
			if iu.Ip == clientIp {
				now := time.Now().Unix()

				if !vdp.IndividualUsers[index].Banned {
					if now < iu.StartTime+int64(vdp.AttackTimespan) {
						vdp.IndividualUsers[index].Count = vdp.IndividualUsers[index].Count + 1
					} else {
						vdp.IndividualUsers[index].Count = 0
						vdp.IndividualUsers[index].StartTime = now
					}
				}
				if ((iu.Count > vdp.AttackCount) ||
					(iu.Count == vdp.AttackCount)) &&
					(now > iu.DateToBeBanned) {
					vdp.IndividualUsers[index].Banned = true

					vdp.BanOccured = true

					vdp.IndividualUsers[index].Count = 0

					AnotherNow := time.Now().Unix() + (int64(vdp.BanTime))

					vdp.IndividualUsers[index].DateToBeBanned = AnotherNow

					vdp.IndividualUsers[index].StartTime = -1

					vdp.BanOccured = true
				}

				repeatedAttack = true

				break
			}
		}
	}

	if !repeatedAttack {
		if len(vdp.Whitelist) > 0 {
			listMemberFound := false

			for _, iu := range vdp.Whitelist {
				if iu == clientIp {
					listMemberFound = true

					break
				}
			}

			if !listMemberFound {
				Client := VDPUser{
					Ip:             clientIp,
					Count:          1,
					StartTime:      time.Now().Unix(),
					Banned:         false,
					DateToBeBanned: -1,
				}

				vdp.IndividualUsers = append(vdp.IndividualUsers, Client)
			}
		} else {
			Client := VDPUser{
				Ip:             clientIp,
				Count:          1,
				StartTime:      time.Now().Unix(),
				Banned:         false,
				DateToBeBanned: -1,
			}

			vdp.IndividualUsers = append(vdp.IndividualUsers, Client)
		}
	}

	return vdp
}

// to-do
func (vdp VanillaDdosProtector) OpenWhiteList(list string) VanillaDdosProtector {
	return vdp
}

func (vdp VanillaDdosProtector) LogEverything() {
	fmt.Printf("attack count: %d\n", vdp.AttackCount)
	fmt.Printf("attack timespan: %d\n", vdp.AttackTimespan)
	fmt.Printf("exact banned time: %d\n", vdp.BanTime)
	fmt.Printf("error code: %d\n", vdp.ErrorCode)
	fmt.Printf("Your users:\n")
	fmt.Println(vdp.IndividualUsers)

}
