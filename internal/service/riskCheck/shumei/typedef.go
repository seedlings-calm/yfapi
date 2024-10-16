package shumei

import "yfapi/core/coreConfig"

const (
	RiskLevelPass   = "PASS"   //通过
	RiskLevelReview = "REVIEW" //可疑
	RiskLevelReject = "REJECT" //违规
)

// 文本检测地址
const (
	TextCheckUrl_BEIJING  = "http://api-text-bj.fengkongcloud.com/text/v4"
	TextCheckUrl_SHANGHAI = "http://api-text-sh.fengkongcloud.com/text/v4"
	TextCheckUrl_MEIGUO   = "http://api-text-fjny.fengkongcloud.com/text/v4"
	TextCheckUrl_XINJIAPO = "http://api-text-xjp.fengkongcloud.com/text/v4"
)

// 图片检测地址
const (
	//批量检测
	ImageCheckUrl_BEIJING  = "http://api-img-bj.fengkongcloud.com/images/v4"
	ImageCheckUrl_SHANGHAI = "http://api-img-sh.fengkongcloud.com/images/v4"
	ImageCheckUrl_MEIGUO   = "http://api-img-gg.fengkongcloud.com/images/v4"
	ImageCheckUrl_XINJIAPO = "http://api-img-xjp.fengkongcloud.com/images/v4"

	//单张检测
	OneImageCheckUrl_BEIJING  = "http://api-img-bj.fengkongcloud.com/image/v4"
	OneImageCheckUrl_SHANGHAI = "http://api-img-sh.fengkongcloud.com/image/v4"
	OneImageCheckUrl_MEIGUO   = "http://api-img-gg.fengkongcloud.com/image/v4"
	OneImageCheckUrl_XINJIAPO = "http://api-img-xjp.fengkongcloud.com/image/v4"
)

// 音频检测地址
const (
	AudioCheckUrl_SHANGHAI = "http://api-audio-sh.fengkongcloud.com/audio/v4"
	AudioCheckUrl_MEIGUO   = "http://api-audio-gg.fengkongcloud.com/audio/v4"
	AudioCheckUrl_XINJIAPO = "http://api-audio-xjp.fengkongcloud.com/audio/v4"

	//同步地址
	AudioCheckUrl = "http://api-audio-sh.fengkongcloud.com/audiomessage/v4"
)

// 视频检测地址
const (
	VideoCheckUrl_BEIJING  = "http://api-video-bj.fengkongcloud.com/video/v4"
	VideoCheckUrl_SHANGHAI = "http://api-video-sh.fengkongcloud.com/video/v4"
	VideoCheckUrl_XINJIAPO = "http://api-video-xjp.fengkongcloud.com/video/v4"
)

type ShuMei struct {
}

type CheckConfig interface {
	getPayLoadData() any
}

func GetVideoUrl() string {
	line := coreConfig.GetHotConf().ShuMei.NetworkLine
	switch line {
	case "SHANGHAI":
		return VideoCheckUrl_SHANGHAI
	case "BEIJING":
		return VideoCheckUrl_BEIJING
	case "XINJIAPO":
		return VideoCheckUrl_XINJIAPO
	default:
		return VideoCheckUrl_BEIJING
	}
}

func GetAudioUrl() string {
	line := coreConfig.GetHotConf().ShuMei.NetworkLine
	switch line {
	case "SHANGHAI":
		return AudioCheckUrl_SHANGHAI
	case "MEIGUO":
		return AudioCheckUrl_MEIGUO
	case "XINJIAPO":
		return AudioCheckUrl_XINJIAPO
	default:
		return AudioCheckUrl_SHANGHAI
	}
}

// 获取图片检测线路地址
func GetImageUrl() string {
	line := coreConfig.GetHotConf().ShuMei.NetworkLine
	switch line {
	case "BEIJING":
		return ImageCheckUrl_BEIJING
	case "SHANGHAI":
		return ImageCheckUrl_SHANGHAI
	case "MEIGUO":
		return ImageCheckUrl_MEIGUO
	case "XINJIAPO":
		return ImageCheckUrl_XINJIAPO
	default:
		return ImageCheckUrl_BEIJING
	}
}

func GetOneImageUrl() string {
	line := coreConfig.GetHotConf().ShuMei.NetworkLine
	switch line {
	case "BEIJING":
		return OneImageCheckUrl_BEIJING
	case "SHANGHAI":
		return OneImageCheckUrl_SHANGHAI
	case "MEIGUO":
		return OneImageCheckUrl_MEIGUO
	case "XINJIAPO":
		return OneImageCheckUrl_XINJIAPO
	default:
		return OneImageCheckUrl_BEIJING
	}
}

func GetTextUrl() string {
	line := coreConfig.GetHotConf().ShuMei.NetworkLine
	switch line {
	case "BEIJING":
		return TextCheckUrl_BEIJING
	case "SHANGHAI":
		return TextCheckUrl_SHANGHAI
	case "MEIGUO":
		return TextCheckUrl_MEIGUO
	case "XINJIAPO":
		return TextCheckUrl_XINJIAPO
	default:
		return TextCheckUrl_BEIJING
	}
}

// 文本检测请求参数
type TextCheckReq struct {
	AccessKey string           `json:"accessKey"`
	AppId     string           `json:"appId"`
	EventId   string           `json:"eventId"`
	Type      string           `json:"type"`
	Data      TextCheckReqData `json:"data"`
}

type TextCheckReqData struct {
	Text     string            `json:"text"`
	TokenId  string            `json:"tokenId"`
	Lang     string            `json:"lang,omitempty"`
	Ip       string            `json:"ip,omitempty"`
	DeviceId string            `json:"deviceId"`
	Nickname string            `json:"nickname"`
	Extra    TextCheckReqExtra `json:"extra"`
}

type TextCheckReqExtra struct {
	Topic          string `json:"topic"`
	AtId           string `json:"atId"`
	Room           string `json:"room"`
	ReceiveTokenId string `json:"receiveTokenId"`
}

// 文本检测返回结构
type TextCheckResp struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	RequestId       string `json:"requestId"`
	FinalResult     int    `json:"finalResult"`
	ResultType      int    `json:"resultType"`
	RiskDescription string `json:"riskDescription"`
	RiskDetail      struct {
	} `json:"riskDetail"`
	RiskLabel1 string `json:"riskLabel1"`
	RiskLabel2 string `json:"riskLabel2"`
	RiskLabel3 string `json:"riskLabel3"`
	RiskLevel  string `json:"riskLevel"`
	AllLabels  []struct {
		Probability     float64 `json:"probability"`
		RiskDescription string  `json:"riskDescription"`
		RiskDetail      struct {
			MatchedLists []struct {
				Name  string `json:"name"`
				Words []struct {
					Position []int  `json:"position"`
					Word     string `json:"word"`
				} `json:"words"`
			} `json:"matchedLists,omitempty"`
		} `json:"riskDetail"`
		RiskLabel1 string `json:"riskLabel1"`
		RiskLabel2 string `json:"riskLabel2"`
		RiskLabel3 string `json:"riskLabel3"`
		RiskLevel  string `json:"riskLevel"`
	} `json:"allLabels"`
	AuxInfo struct {
		ContactResult []struct {
			ContactString string `json:"contactString"`
			ContactType   int    `json:"contactType"`
		} `json:"contactResult"`
		FilteredText string `json:"filteredText"`
	} `json:"auxInfo"`
	BusinessLabels []interface{} `json:"businessLabels"`
}

// 图片异步检测
type ImagesAsyncCheckReq struct {
	AccessKey string                  `json:"accessKey"`
	AppId     string                  `json:"appId"`
	Callback  string                  `json:"callback"`
	Data      ImagesAsyncCheckReqData `json:"data"`
	EventId   string                  `json:"eventId"`
	Type      string                  `json:"type"`
}

type ImagesAsyncCheckReqData struct {
	Imgs    []ImagesAsyncCheckReqImgs `json:"imgs"`
	TokenId string                    `json:"tokenId"`
}

type ImagesAsyncCheckReqImgs struct {
	BtId string `json:"btId"`
	Img  string `json:"img"`
}

// 图片异步检测返回结果
type ImagesAsyncCheckResp struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	RequestIds []struct {
		BtId      string `json:"btId"`
		RequestId string `json:"requestId"`
	} `json:"requestIds"`
}

// 图片异步检测回调结果
type ImagesCheckCallBackResp struct {
	AuxInfo struct {
	} `json:"auxInfo"`
	Code int `json:"code"`
	Imgs []struct {
		AllLabels []struct {
			Probability     float64 `json:"probability"`
			RiskDescription string  `json:"riskDescription"`
			RiskDetail      struct {
				OcrText struct {
					Text string `json:"text"`
				} `json:"ocrText"`
				RiskSource int `json:"riskSource"`
			} `json:"riskDetail"`
			RiskLabel1 string `json:"riskLabel1"`
			RiskLabel2 string `json:"riskLabel2"`
			RiskLabel3 string `json:"riskLabel3"`
			RiskLevel  string `json:"riskLevel"`
		} `json:"allLabels"`
		AuxInfo struct {
			Segments    int `json:"segments"`
			TypeVersion struct {
				BAN      string `json:"BAN,omitempty"`
				MINOR    string `json:"MINOR,omitempty"`
				OCR      string `json:"OCR"`
				POLITICS string `json:"POLITICS,omitempty"`
				PORN     string `json:"PORN"`
				VIOLENCE string `json:"VIOLENCE,omitempty"`
			} `json:"typeVersion"`
		} `json:"auxInfo"`
		BusinessLabels []struct {
			BusinessDescription string `json:"businessDescription"`
			BusinessDetail      struct {
				FaceRatio float64 `json:"face_ratio,omitempty"`
				Faces     []struct {
					FaceRatio   float64 `json:"face_ratio"`
					Id          string  `json:"id"`
					Location    []int   `json:"location"`
					Name        string  `json:"name"`
					Probability float64 `json:"probability"`
				} `json:"faces,omitempty"`
				Location    []int   `json:"location,omitempty"`
				Name        string  `json:"name,omitempty"`
				Probability float64 `json:"probability,omitempty"`
				FaceNum     int     `json:"face_num,omitempty"`
			} `json:"businessDetail"`
			BusinessLabel1  string  `json:"businessLabel1"`
			BusinessLabel2  string  `json:"businessLabel2"`
			BusinessLabel3  string  `json:"businessLabel3"`
			ConfidenceLevel int     `json:"confidenceLevel"`
			Probability     float64 `json:"probability"`
		} `json:"businessLabels,omitempty"`
		Code            int    `json:"code"`
		FinalResult     int    `json:"finalResult,omitempty"`
		Message         string `json:"message"`
		RequestId       string `json:"requestId"`
		ResultType      int    `json:"resultType,omitempty"`
		RiskDescription string `json:"riskDescription"`
		RiskDetail      struct {
			OcrText struct {
				Text string `json:"text"`
			} `json:"ocrText,omitempty"`
			RiskSource int `json:"riskSource"`
		} `json:"riskDetail"`
		RiskLabel1  string `json:"riskLabel1"`
		RiskLabel2  string `json:"riskLabel2"`
		RiskLabel3  string `json:"riskLabel3"`
		RiskLevel   string `json:"riskLevel"`
		RiskSource  int    `json:"riskSource,omitempty"`
		TokenLabels struct {
			UGCAccountRisk struct {
			} `json:"UGC_account_risk"`
		} `json:"tokenLabels"`
		BtId string `json:"btId,omitempty"`
	} `json:"imgs"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

// 音频检测请求结构
type AudioCheckReq struct {
	AccessKey    string            `json:"accessKey"`
	AppId        string            `json:"appId"`
	EventId      string            `json:"eventId"`
	Type         string            `json:"type"`
	BusinessType string            `json:"businessType"`
	BtId         string            `json:"btId"`
	ContentType  string            `json:"contentType"`
	Content      string            `json:"content"`
	Callback     string            `json:"callback"`
	Data         AudioCheckReqData `json:"data"`
}

type AudioCheckReqData struct {
	ReturnAllText int    `json:"returnAllText"`
	TokenId       string `json:"tokenId"`
}

// 返回结果
type AudioCheckResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	BtId      string `json:"btId"`
}

// 音频同步返回结果
type AudioCheckSyncResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	BtId      string `json:"btId"`
	Detail    struct {
		AudioDetail []struct {
			RequestId      string `json:"requestId"`
			AudioStarttime int    `json:"audioStarttime"`
			AudioEndtime   int    `json:"audioEndtime"`
			AudioUrl       string `json:"audioUrl"`
			BusinessLabels []struct {
				BusinessDescription string `json:"businessDescription"`
				BusinessDetail      struct {
				} `json:"businessDetail"`
				BusinessLabel1  string  `json:"businessLabel1"`
				BusinessLabel2  string  `json:"businessLabel2"`
				BusinessLabel3  string  `json:"businessLabel3"`
				ConfidenceLevel int     `json:"confidenceLevel"`
				Probability     float64 `json:"probability"`
			} `json:"businessLabels"`
			AllLabels       []interface{} `json:"allLabels"`
			RiskLevel       string        `json:"riskLevel"`
			RiskLabel1      string        `json:"riskLabel1"`
			RiskLabel2      string        `json:"riskLabel2"`
			RiskLabel3      string        `json:"riskLabel3"`
			RiskDescription string        `json:"riskDescription"`
			RiskDetail      struct {
				AudioText string `json:"audioText"`
			} `json:"riskDetail"`
		} `json:"audioDetail"`
		AudioTags struct {
			Gender struct {
				Label       string `json:"label"`
				Probability int    `json:"probability"`
			} `json:"gender"`
			Language []struct {
				Confidence int `json:"confidence"`
				Label      int `json:"label"`
			} `json:"language"`
			Song   int `json:"song"`
			Timbre []struct {
				Label       string `json:"label"`
				Probability int    `json:"probability"`
			} `json:"timbre"`
		} `json:"audioTags"`
		AudioText     string `json:"audioText"`
		AudioTime     int    `json:"audioTime"`
		Code          int    `json:"code"`
		RequestParams struct {
			Channel       string `json:"channel"`
			Lang          string `json:"lang"`
			ReturnAllText int    `json:"returnAllText"`
			TokenId       string `json:"tokenId"`
		} `json:"requestParams"`
		RiskLevel string `json:"riskLevel"`
	} `json:"detail"`
}

// 音频检测回调结果
type AudioCheckCallbackResult struct {
	RequestId   string `json:"requestId"`
	BtId        string `json:"btId"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RiskLevel   string `json:"riskLevel"`
	AudioDetail []struct {
		RequestId      string `json:"requestId"`
		AudioStarttime int    `json:"audioStarttime"`
		AudioEndtime   int    `json:"audioEndtime"`
		AudioUrl       string `json:"audioUrl"`
		BusinessLabels []struct {
			BusinessDescription string `json:"businessDescription"`
			BusinessDetail      struct {
			} `json:"businessDetail"`
			BusinessLabel1  string  `json:"businessLabel1"`
			BusinessLabel2  string  `json:"businessLabel2"`
			BusinessLabel3  string  `json:"businessLabel3"`
			ConfidenceLevel int     `json:"confidenceLevel"`
			Probability     float64 `json:"probability"`
		} `json:"businessLabels,omitempty"`
		AllLabels []struct {
			Probability     int    `json:"probability"`
			RiskDescription string `json:"riskDescription"`
			RiskDetail      struct {
				RiskSource int `json:"riskSource"`
			} `json:"riskDetail"`
			RiskLabel1 string `json:"riskLabel1"`
			RiskLabel2 string `json:"riskLabel2"`
			RiskLabel3 string `json:"riskLabel3"`
			RiskLevel  string `json:"riskLevel"`
		} `json:"allLabels,omitempty"`
		RiskLevel       string `json:"riskLevel"`
		RiskLabel1      string `json:"riskLabel1"`
		RiskLabel2      string `json:"riskLabel2"`
		RiskLabel3      string `json:"riskLabel3"`
		RiskDescription string `json:"riskDescription"`
		RiskDetail      struct {
			AudioText string `json:"audioText"`
		} `json:"riskDetail,omitempty"`
	} `json:"audioDetail"`
	AudioTags struct {
		Gender struct {
			Label       string `json:"label"`
			Probability int    `json:"probability"`
		} `json:"gender"`
		Language []struct {
			Confidence int `json:"confidence"`
			Label      int `json:"label"`
		} `json:"language"`
		Song   int `json:"song"`
		Timbre []struct {
			Label       string `json:"label"`
			Probability int    `json:"probability"`
		} `json:"timbre"`
	} `json:"audioTags"`
}

// 单张图片同步检测请求参数
type OneImageSyncCheckReq struct {
	AccessKey string                   `json:"accessKey"`
	AppId     string                   `json:"appId"`
	Data      OneImageSyncCheckReqData `json:"data"`
	EventId   string                   `json:"eventId"`
	Type      string                   `json:"type"`
}

type OneImageSyncCheckReqData struct {
	Img     string `json:"img"`
	TokenId string `json:"tokenId"`
}

// 单张图片同步检测返回结果
type OneImageSyncCheckResp struct {
	AllLabels []struct {
		Probability     float64 `json:"probability"`
		RiskDescription string  `json:"riskDescription"`
		RiskDetail      struct {
			OcrText struct {
				Text string `json:"text"`
			} `json:"ocrText"`
			RiskSource int `json:"riskSource"`
		} `json:"riskDetail"`
		RiskLabel1 string `json:"riskLabel1"`
		RiskLabel2 string `json:"riskLabel2"`
		RiskLabel3 string `json:"riskLabel3"`
		RiskLevel  string `json:"riskLevel"`
	} `json:"allLabels"`
	AuxInfo struct {
		Segments    int `json:"segments"`
		TypeVersion struct {
			BAN      string `json:"BAN"`
			MINOR    string `json:"MINOR"`
			OCR      string `json:"OCR"`
			POLITICS string `json:"POLITICS"`
			PORN     string `json:"PORN"`
			VIOLENCE string `json:"VIOLENCE"`
		} `json:"typeVersion"`
	} `json:"auxInfo"`
	BusinessLabels []struct {
		BusinessDescription string `json:"businessDescription"`
		BusinessDetail      struct {
			FaceRatio float64 `json:"face_ratio,omitempty"`
			Faces     []struct {
				FaceRatio   float64 `json:"face_ratio"`
				Id          string  `json:"id"`
				Location    []int   `json:"location"`
				Name        string  `json:"name"`
				Probability float64 `json:"probability"`
			} `json:"faces,omitempty"`
			Location    []int   `json:"location,omitempty"`
			Name        string  `json:"name,omitempty"`
			Probability float64 `json:"probability,omitempty"`
			FaceNum     int     `json:"face_num,omitempty"`
		} `json:"businessDetail"`
		BusinessLabel1  string  `json:"businessLabel1"`
		BusinessLabel2  string  `json:"businessLabel2"`
		BusinessLabel3  string  `json:"businessLabel3"`
		ConfidenceLevel int     `json:"confidenceLevel"`
		Probability     float64 `json:"probability"`
	} `json:"businessLabels"`
	Code            int    `json:"code"`
	FinalResult     int    `json:"finalResult"`
	Message         string `json:"message"`
	RequestId       string `json:"requestId"`
	ResultType      int    `json:"resultType"`
	RiskDescription string `json:"riskDescription"`
	RiskDetail      struct {
		OcrText struct {
			Text string `json:"text"`
		} `json:"ocrText"`
		RiskSource int `json:"riskSource"`
	} `json:"riskDetail"`
	RiskLabel1  string `json:"riskLabel1"`
	RiskLabel2  string `json:"riskLabel2"`
	RiskLabel3  string `json:"riskLabel3"`
	RiskLevel   string `json:"riskLevel"`
	RiskSource  int    `json:"riskSource"`
	TokenLabels struct {
		UGCAccountRisk struct {
		} `json:"UGC_account_risk"`
	} `json:"tokenLabels"`
}

type VideoCheckReq struct {
	AccessKey       string            `json:"accessKey"`
	AppId           string            `json:"appId"`
	AudioType       string            `json:"audioType"`
	Callback        string            `json:"callback"`
	Data            VideoCheckReqData `json:"data"`
	EventId         string            `json:"eventId"`
	ImgBusinessType string            `json:"imgBusinessType"`
	ImgType         string            `json:"imgType"`
}

type VideoCheckReqData struct {
	BtId    string `json:"btId"`
	TokenId string `json:"tokenId"`
	Url     string `json:"url"`
}

type VideoCheckResp struct {
	BtId      string `json:"btId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

type VideoCheckCallbackResult struct {
	AudioDetail []struct {
		AllLabels []struct {
			Probability     float64 `json:"probability"`
			RiskDescription string  `json:"riskDescription"`
			RiskDetail      struct {
				AudioText  string `json:"audioText"`
				RiskSource int    `json:"riskSource"`
			} `json:"riskDetail"`
			RiskLabel1 string `json:"riskLabel1"`
			RiskLabel2 string `json:"riskLabel2"`
			RiskLabel3 string `json:"riskLabel3"`
			RiskLevel  string `json:"riskLevel"`
		} `json:"allLabels"`
		AudioEndtime    int           `json:"audioEndtime"`
		AudioStarttime  int           `json:"audioStarttime"`
		AudioText       string        `json:"audioText"`
		AudioUrl        string        `json:"audioUrl"`
		BusinessLabels  []interface{} `json:"businessLabels"`
		RequestId       string        `json:"requestId"`
		RiskDescription string        `json:"riskDescription"`
		RiskDetail      struct {
			AudioText  string `json:"audioText"`
			RiskSource int    `json:"riskSource"`
		} `json:"riskDetail"`
		RiskLabel1 string `json:"riskLabel1"`
		RiskLabel2 string `json:"riskLabel2"`
		RiskLabel3 string `json:"riskLabel3"`
		RiskLevel  string `json:"riskLevel"`
	} `json:"audioDetail"`
	AuxInfo struct {
		BillingAudioDuration int `json:"billingAudioDuration"`
		BillingImgNum        int `json:"billingImgNum"`
		FrameCount           int `json:"frameCount"`
		Time                 int `json:"time"`
	} `json:"auxInfo"`
	BtId        string `json:"btId"`
	Code        int    `json:"code"`
	FrameDetail []struct {
		AllLabels []struct {
			Probability     float64 `json:"probability"`
			RiskDescription string  `json:"riskDescription"`
			RiskDetail      struct {
				OcrText struct {
					Text string `json:"text"`
				} `json:"ocrText"`
				RiskSource int `json:"riskSource"`
			} `json:"riskDetail"`
			RiskLabel1 string `json:"riskLabel1"`
			RiskLabel2 string `json:"riskLabel2"`
			RiskLabel3 string `json:"riskLabel3"`
			RiskLevel  string `json:"riskLevel"`
		} `json:"allLabels"`
		AuxInfo struct {
			Similarity float64 `json:"similarity"`
		} `json:"auxInfo"`
		BusinessLabels []struct {
			BusinessDescription string `json:"businessDescription"`
			BusinessDetail      struct {
				FaceNum int `json:"face_num,omitempty"`
				Faces   []struct {
					FaceRatio   float64 `json:"face_ratio"`
					Id          string  `json:"id"`
					Location    []int   `json:"location"`
					Name        string  `json:"name"`
					Probability float64 `json:"probability"`
				} `json:"faces,omitempty"`
			} `json:"businessDetail"`
			BusinessLabel1  string  `json:"businessLabel1"`
			BusinessLabel2  string  `json:"businessLabel2"`
			BusinessLabel3  string  `json:"businessLabel3"`
			ConfidenceLevel int     `json:"confidenceLevel"`
			Probability     float64 `json:"probability"`
		} `json:"businessLabels"`
		ImgText         string `json:"imgText"`
		ImgUrl          string `json:"imgUrl"`
		RequestId       string `json:"requestId"`
		RiskDescription string `json:"riskDescription"`
		RiskDetail      struct {
			OcrText struct {
				Text string `json:"text"`
			} `json:"ocrText"`
			RiskSource int `json:"riskSource"`
		} `json:"riskDetail"`
		RiskLabel1 string `json:"riskLabel1"`
		RiskLabel2 string `json:"riskLabel2"`
		RiskLabel3 string `json:"riskLabel3"`
		RiskLevel  string `json:"riskLevel"`
		Time       int    `json:"time"`
	} `json:"frameDetail"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	RiskLevel string `json:"riskLevel"`
}
