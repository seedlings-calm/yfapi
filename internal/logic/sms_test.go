package logic

import (
	"context"
	"testing"
)

func TestSendMsg(t *testing.T) {
	var req = []Sms{
		{RegionCode: "+86", Mobile: "13103810741", Type: 1, SignType: "sha1"},
		{RegionCode: "+86", Mobile: "1310381", Type: 1, SignType: "sha1"},
		{RegionCode: "+86", Mobile: "18516179582", Type: 2, SignType: "md5"},
		{RegionCode: "", Mobile: "13103810741", Type: 2, SignType: ""},
		{RegionCode: "+86", Mobile: "1310381", Type: 3, SignType: ""},
	}
	for _, v := range req {
		err := v.SendMsg()
		t.Log(err)

	}

}

func TestCheckSms(t *testing.T) {
	var req = []Sms{
		{RegionCode: "+86", Mobile: "13103810741", Type: 1, SignType: "sha1"},
		{RegionCode: "+86", Mobile: "1310381", Type: 1, SignType: "sha1"},
		{RegionCode: "+86", Mobile: "18516179582", Type: 2, SignType: "md5"},
		{RegionCode: "", Mobile: "13103810741", Type: 2, SignType: ""},
		{RegionCode: "+86", Mobile: "1310381", Type: 3, SignType: ""},
	}
	for _, v := range req {
		err := v.CheckSms(context.Background())
		t.Log(err)
	}
}
