package image

import "testing"

func TestImageUploader_Upload(t *testing.T) {
	cookie := "SINAGLOBAL=2277809656598.042.1507216888871; _s_tentry=-; YF-Page-G0=9a31b867b34a0b4839fa27a4ab6ec79f; Apache=5351375581362.308.1524565902478; ULV=1524565902942:11:1:1:5351375581362.308.1524565902478:1521951597186; YF-Ugrow-G0=ad83bc19c1269e709f753b172bddb094; YF-V5-G0=7fb6f47dfff7c4352ece66bba44a6e5a; login_sid_t=e2a015a3417817d7ede7236e699a75fc; cross_origin_proto=SSL; SSOLoginState=1526979421; UOR=f.uliba.net,widget.weibo.com,www.imydl.com; wb_view_log=1920*12001; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WFKIokbz3g-4Tq3R53_.Drn5JpX5K2hUgL.Foe4S054eo.N1h52dJLoIp7LxKML1KBLBKnLxKqL1hnLBoUfUgvrIP5f9rya; ALF=1560129373; SCF=AimoCxlxenffPp5foS0MfzEW3JjBxiTCa3vmnDkOSNRmXytfUhI20nwS38AwxHb_I8B3zA1-UMz8l673rxIsHlI.; SUB=_2A252GAuPDeRhGeVH7FIY8ifLwzyIHXVVbHpHrDV8PUNbmtBeLVHBkW9NTyFfFl5Fb0dDenC7RZk7_wfThTUj3wEW; SUHB=079jaewNo9yM5S; un=18701381686; wvr=6"
	uploader := NewUploader(cookie)
	uploader.Upload("./2.jpg")
}
