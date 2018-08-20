package event

import (
	"testing"
	"fmt"
	"unsafe"
	"reflect"
)

type EventHandler func(event *Event)

func eventHandler2(event *Event) {
	fmt.Print("event is %+v", event)
}
func eventHandler3(event *Event){
	fmt.Print("event is %+v", event)
}
func Test_Test(t *testing.T) {

	h := eventHandler2
	s := []byte(str)
	p := *(*EventHandler)(unsafe.Pointer(&h))
	s = []byte(str)
	x := *(*EventHandler)(unsafe.Pointer(&h))
	s = []byte(str)
	//faked := *(*func(float64))(unsafe.Pointer(&fn))
	a := uintptr(unsafe.Pointer(&h))
	s = []byte(str)
	b := uintptr(unsafe.Pointer(&x))
	s = []byte(str)
	fmt.Print(p, a,b, h, eventHandler2, eventHandler3, s[0])


}

func Test_Test2(t *testing.T) {
	a := eventHandler2
	//b := eventHandler2
	e := eventHandler3
	c := reflect.ValueOf(a)
	d := reflect.ValueOf(e)
	fmt.Println(c == d)
}


var str = `{"graphql":{"shortcode_media":{"__typename":"GraphVideo","id":"1747067984893244860","shortcode":"Bg-1RNuBPW8","dimensions":{"height":480,"width":480},"gating_info":null,"media_preview":"ACoqm02PCM57nH+fxp17IyISDg9v6/jVi3GyJR3xn8+ahFuJpGMp+XACgH8zReysVZvUz7OF51diWVwAUbJwfb3q/aTF/kcYdev86mtI5IwQ7qIkzgDqQfUnpjt3zTYk8stlt+9tyt3xjHJHWkmJoh1BN0W7upz+HesX8q6WRd6keoxXO5I49PpRIDdkUheThR6Uli6vIQvTb1/Gq1zcbgVHNVLK8WByE+Yk4/AelQl1NXLobLwmRzHxtYfNUDwvA68gx42jHGD7/X+dWUuhycdarmU3AOFb5WHUYHH1xnHtVIzZZzWU9luYn1JNaKtxxTM1ZJmSJhSo6+tZEY2uCK2m6Vjyff8AxpDNQX/2dsBdx963lLCMMfTJHuea5KfrXRXTEW2c/wAP9KljK1vP5o3YwMkD8Kt1mad0X3/xrRrQk//Z","display_url":"https://scontent-lax3-2.cdninstagram.com/vp/9dd3923db62fe2e712cba50950f95d03/5AC1AE51/t51.2885-15/e35/29417734_411935732585496_1353847354683293696_n.jpg","display_resources":[{"src":"https://scontent-lax3-2.cdninstagram.com/vp/9dd3923db62fe2e712cba50950f95d03/5AC1AE51/t51.2885-15/e35/29417734_411935732585496_1353847354683293696_n.jpg","config_width":640,"config_height":640},{"src":"https://scontent-lax3-2.cdninstagram.com/vp/9dd3923db62fe2e712cba50950f95d03/5AC1AE51/t51.2885-15/e35/29417734_411935732585496_1353847354683293696_n.jpg","config_width":750,"config_height":750},{"src":"https://scontent-lax3-2.cdninstagram.com/vp/9dd3923db62fe2e712cba50950f95d03/5AC1AE51/t51.2885-15/e35/29417734_411935732585496_1353847354683293696_n.jpg","config_width":1080,"config_height":1080}],"dash_info":{"is_dash_eligible":false,"video_dash_manifest":null,"number_of_qualities":0},"video_url":"https://scontent-lax3-2.cdninstagram.com/vp/8a079411f1caeaf2a04b2aaa7d8d84a7/5AC1D49B/t50.2886-16/29912446_1753695857986743_13364288227598477_n.mp4","video_view_count":32,"is_video":true,"should_log_client_event":false,"tracking_token":"eyJ2ZXJzaW9uIjo1LCJwYXlsb2FkIjp7ImlzX2FuYWx5dGljc190cmFja2VkIjpmYWxzZSwidXVpZCI6IjljYjk3NzcyZTQ2NjRkNDQ5MDc4MzQ5ZjhhOGY2OTg0MTc0NzA2Nzk4NDg5MzI0NDg2MCJ9LCJzaWduYXR1cmUiOiIifQ==","edge_media_to_tagged_user":{"edges":[]},"edge_media_to_caption":{"edges":[{"node":{"text":"The world of all cute animals,\nshare your pictures with us,\nplease tag us for feature.\nby  @mywinterfells.siberian_huskies\n.\n.\n.\n.\n#dogwalker #lovemydog #animalphotography #petlovers #petlover #petsagram #animalprint #petsofig #rescuedog #cutepetclub #instadogs #petsofinstagram #petshop #animalkingdom #bulldogsofinstagram #like4like #cuteguy #puppiesofinstagram #pitbullsofinstagram #cuteness #toocute #dogsofinsta #ilovedogs #cutiepatootie #tbt #cutenessoverload #chillin #dogstagram #dogmodel #instacat"}}]},"caption_is_edited":false,"edge_media_to_comment":{"count":1,"page_info":{"has_next_page":false,"end_cursor":null},"edges":[{"node":{"id":"17906410561156323","text":"@happy.quokka @_saraah_xoxo","created_at":1522488497,"owner":{"id":"242965203","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/a47975ce1d10f10323c60bed7c8dcca2/5B55E998/t51.2885-19/s150x150/29093119_182191412568989_4192371732868235264_n.jpg","username":"santipanti88"}}}]},"comments_disabled":false,"taken_at_timestamp":1522486771,"edge_media_preview_like":{"count":23,"edges":[{"node":{"id":"6335679301","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/630a28390a6bf496c595063c23c47a66/5B51DF9B/t51.2885-19/s150x150/29089518_2123074787706854_2549707541428830208_n.jpg","username":"skythegreatdanex"}},{"node":{"id":"4187783858","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/a7031b0f919ef186e643c0c0c1bdb13e/5B4FE782/t51.2885-19/s150x150/29093311_190543494895362_2550172067911696384_n.jpg","username":"jus_an0ther.bean"}},{"node":{"id":"2873799046","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/5a25e78c71602b7a1a615683b6acc57e/5B5CF7B7/t51.2885-19/s150x150/18096129_1759249931051766_1719925413234343936_a.jpg","username":"georgethetoypoodle"}},{"node":{"id":"357008551","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/30cac5d753953ac4b51ef83d1e9f975e/5B74002E/t51.2885-19/s150x150/26871791_400657853703386_5067201689971326976_n.jpg","username":"miriamfernandezz__"}},{"node":{"id":"5447583017","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/5c7e2ba20cf2c39fec7d3dec115ed740/5B752F4D/t51.2885-19/s150x150/26071337_1387669334676942_8507813747853623296_n.jpg","username":"rosie_the_cavalierx"}},{"node":{"id":"1351927247","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/06070257fcb50efeaf38d036a1dd1fd1/5B535062/t51.2885-19/s150x150/13774180_580281578804294_1247522296_a.jpg","username":"paul_dickel"}},{"node":{"id":"1703056972","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/25a3bb0e9a67027977b3150402627d5a/5B7578E5/t51.2885-19/s150x150/28155709_549855178733307_2239775880941404160_n.jpg","username":"laura__sanz"}},{"node":{"id":"7077336992","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/ad79ea233713baac3cd988b42ea3d395/5B72F3E4/t51.2885-19/s150x150/28751372_974819809334894_1598752170927194112_n.jpg","username":"anna._.0102"}},{"node":{"id":"612068573","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/35c4087239177463a6435ccbdd8e5032/5B6A0824/t51.2885-19/s150x150/20184144_1504211752968871_5628247976448098304_a.jpg","username":"tami_kasebacher"}},{"node":{"id":"4857510155","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/083ace402e7026a2f8b1503357570aad/5B6993C4/t51.2885-19/s150x150/19429451_325631081184097_3213927107088351232_a.jpg","username":"raghuveer787"}}]},"edge_media_to_sponsor_user":{"edges":[]},"location":null,"viewer_has_liked":false,"viewer_has_saved":false,"viewer_has_saved_to_collection":false,"owner":{"id":"498059394","profile_pic_url":"https://scontent-lax3-2.cdninstagram.com/vp/02c48406f55cfce5ef6cb2b15c81a478/5B7599A7/t51.2885-19/s150x150/26070371_1602111373203006_931916398953758720_n.jpg","username":"cuty__animal_","blocked_by_viewer":false,"followed_by_viewer":false,"full_name":"All cute animals","has_blocked_viewer":false,"is_private":false,"is_unpublished":false,"is_verified":false,"requested_by_viewer":false},"is_ad":false,"edge_web_media_to_related_media":{"edges":[]}}}}`